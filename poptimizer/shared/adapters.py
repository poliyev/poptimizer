"""Базовые классы взаимодействия с внешней инфраструктурой."""
import asyncio
import logging
import weakref
from collections.abc import MutableMapping
from typing import Any, Callable, Final, Generic, NamedTuple, Optional, TypeVar

from motor import motor_asyncio
from pymongo.collection import Collection

from poptimizer.shared import connections, domain

# Коллекция для сохранения объектов из групп с одним объектом
MISC: Final = "misc"


class AsyncLogger:
    """Асинхронное логирование в отдельном потоке.

    Поддерживает протокол дескриптора для автоматического определения имени класса, в котором он
    является атрибутом.
    """

    def __init__(self) -> None:
        """Инициализация логгера."""
        self._logger = logging.getLogger()

    def __call__(self, message: str) -> None:
        """Создает асинхронную задачу по логгированию."""
        asyncio.create_task(self._logging_task(message))

    def __set_name__(self, owner: type[object], name: str) -> None:
        """Создает логгер с именем класса, где он является атрибутом."""
        self._logger = logging.getLogger(owner.__name__)

    def __get__(self, instance: object, owner: type[object]) -> "AsyncLogger":
        """Возвращает себя при обращении к атрибуту."""
        return self

    async def _logging_task(self, message: str) -> None:
        """Задание по логгированию."""
        await asyncio.to_thread(self._logger.info, message)


class Desc(NamedTuple):
    """Описание кодирования и декодирования из документа MongoDB."""

    field_name: str
    doc_name: str
    factory_name: str
    encoder: Optional[Callable[[Any], Any]] = None  # type: ignore
    decoder: Optional[Callable[[Any], Any]] = None  # type: ignore


EntityType = TypeVar("EntityType", bound=domain.BaseEntity)


class Mapper(Generic[EntityType]):
    """Сохраняет и загружает доменные объекты из MongoDB."""

    _identity_map: MutableMapping[
        domain.ID,
        EntityType,
    ] = weakref.WeakValueDictionary()
    _logger = AsyncLogger()

    def __init__(  # type: ignore
        self,
        desc_list: tuple[Desc, ...],
        factory: domain.AbstractFactory[EntityType],
        client: motor_asyncio.AsyncIOMotorClient = connections.MONGO_CLIENT,
    ) -> None:
        """Сохраняет соединение с MongoDB, информацию для мэппинга объектов и фабрику."""
        self._client = client
        self._desc_list = desc_list
        self._factory = factory

    async def __call__(self, id_: domain.ID) -> EntityType:
        """Загружает доменный объект из базы."""
        if (table_old := self._identity_map.get(id_)) is not None:
            return table_old

        mongo_dict = await self.get_doc(id_)
        table = self._decode(id_, mongo_dict)

        if (table_old := self._identity_map.get(id_)) is not None:
            return table_old

        self._identity_map[id_] = table

        return table

    async def get_doc(self, id_: domain.ID) -> domain.StateDict:
        """Запрашивает документ по ID.

        При отсутствии возвращает пустой словарь.
        """
        collection, name = self._get_collection_and_id(id_)
        return await collection.find_one({"_id": name}, projection={"_id": False}) or {}

    async def commit(
        self,
        entity: EntityType,
    ) -> None:
        """Записывает изменения доменного объекта в MongoDB."""
        if mongo_dict := self._encode(entity):
            id_ = entity.id_
            self._logger(f"Сохранение {id_}")
            collection, name = self._get_collection_and_id(id_)
            await collection.replace_one(
                filter={"_id": name},
                replacement=dict(_id=name, **mongo_dict),
                upsert=True,
            )

    def _get_collection_and_id(self, id_: domain.ID) -> tuple[Collection, str]:
        """Коллекцию и ID документа.

        При совпадении названия группы и имени выбирает специальную коллекцию для одиночных записей.
        """
        collection = id_.group
        name = id_.name
        if collection == name:
            collection = MISC
        return self._client[id_.package][collection], name

    def _encode(self, entity: EntityType) -> domain.StateDict:
        """Кодирует данные в совместимый с MongoDB формат."""
        if not (entity_state := entity.changed_state()):
            return {}

        entity.clear()
        sentinel = object()
        for desc in self._desc_list:
            if (field_value := entity_state.pop(desc.field_name, sentinel)) is sentinel:
                continue
            if desc.encoder:
                field_value = desc.encoder(field_value)
            entity_state[desc.doc_name] = field_value

        return entity_state

    def _decode(self, id_: domain.ID, mongo_dict: domain.StateDict) -> EntityType:
        """Декодирует данные из формата MongoDB формат атрибутов модели и создает объект."""
        sentinel = object()
        for desc in self._desc_list:
            if (field_value := mongo_dict.pop(desc.doc_name, sentinel)) is sentinel:
                continue
            if desc.decoder:
                field_value = desc.decoder(field_value)
            mongo_dict[desc.factory_name] = field_value
        return self._factory(id_, mongo_dict)
