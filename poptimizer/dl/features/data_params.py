"""Описание модели и данных."""
import abc
import copy
from typing import Generator, Tuple

import pandas as pd

from poptimizer import config
from poptimizer.data.views import quotes

FORECAST_DAYS = config.FORECAST_DAYS


def div_price_train_size(
    tickers: Tuple[str, ...],
    end: pd.Timestamp,
) -> Tuple[pd.DataFrame, pd.DataFrame, int]:
    """Данные по дивидендам, ценам и количество дней в тренировочном наборе."""
    div, price = quotes.div_and_prices(tickers, end)
    train_size = len(price) - FORECAST_DAYS

    return div, price, train_size


class DataParams(abc.ABC):
    """Параметры данных для DL-модели."""

    def __init__(self, tickers: Tuple[str, ...], end: pd.Timestamp, params: dict):
        """Модель строится для определенного набора тикеров и диапазона дат.

        Наборы данных для обучения, валидации, тестирования и прогнозирования реализуются в конкретных
        классах. Кроме собственно необходимых для построения признаков параметров класс хранит
        кешированные и обрезанные у четом типа данных и отсутствующих значений информацию о дивидендах
        и стоимости акций, которые могут быть использованы для построения признаков и корректного их
        выравнивания по времени.

        :param tickers:
            Перечень тикеров, для которых будет строится модель.
        :param end:
            Конец диапазона дат статистики по ценам и дивидендам, которые будут использоваться для
            построения модели.
        :param params:
            Словарь с параметрами для построения признаков и других элементов модели.
        """
        self._cache = {}
        self._tickers = tickers
        self._end = end
        self._params = copy.deepcopy(params)
        div, price = self._div_price(tickers, end)
        self._div = {}
        self._price = {}
        for ticker in tickers:
            start = price[ticker].first_valid_index()
            self._div[ticker] = div.loc[start:, ticker]
            self._price[ticker] = price.loc[start:, ticker]

    @property
    def cache(self) -> dict:
        """Словарь для кеширования.

        Признак или схожие признаки могут сохранять вспомогательную информацию, чтобы исключить
        повторные вызовы тяжелых функций.
        """
        return self._cache

    @property
    def tickers(self) -> Tuple[str]:
        """Перечень тикеров."""
        return self._tickers

    @property
    def end(self) -> pd.Timestamp:
        """Конечная дата."""
        return self._end

    @property
    def shuffle(self) -> bool:
        """Нужно ли перемешивать данные."""
        return False

    @property
    def history_days(self) -> int:
        """Длинна истории для признаков в днях."""
        return self._params["history_days"]

    @property
    def batch_size(self) -> int:
        """Размер батча."""
        return self._params["batch_size"]

    def price(self, ticker: str) -> pd.Series:
        """Цены для тикера.

        Обрезаются с учетом диапазона дат, который будет использоваться для построения признаков, с учетом возможного
        отсутствия котировок в начале.
        """
        return self._price[ticker]

    def div(self, ticker: str) -> pd.Series:
        """Дивиденды для тикера.

        Обрезаются с учетом диапазона дат, который будет использоваться для построения признаков, с учетом возможного
        отсутствия котировок в начале.
        """
        return self._div[ticker]

    def len(self, ticker) -> int:
        """Количество доступных примеров для данного тикера."""
        return max(0, len(self.price(ticker)) - self.history_days - FORECAST_DAYS + 1)

    def get_all_feat(self) -> Generator[str, None, None]:
        """Получить все названия признаков."""
        yield from (feat for feat, param in self._params["features"].items() if param["on"])

    def get_feat_params(self, feat_name: str) -> dict:
        """Получить параметры для признака."""
        return self._params["features"][feat_name]

    @abc.abstractmethod
    def _div_price(self, tickers, end) -> Tuple[pd.DataFrame, pd.DataFrame]:
        """Данные о дивидендах и стоимости акций для указанных тикеров и конечной даты.

        Метод должен реализовывать необходимую обрезку с учетом конкретного класса - для обучения,
        валидации, тестирования или тренировки. При необходимости изменять параметры.
        """


class TrainParams(DataParams):
    """Используются признаки, как есть для начальной части семпла."""

    @property
    def shuffle(self):
        """Нужно перемешивать данные."""
        return True

    def _div_price(self, tickers, end) -> Tuple[pd.DataFrame, pd.DataFrame]:
        div, price, train_size = div_price_train_size(tickers, end)
        div = div.iloc[:train_size]
        price = price.iloc[:train_size]

        return div, price


class TestParams(DataParams):
    """Параметры для тестирования."""

    def _div_price(self, tickers, end) -> Tuple[pd.DataFrame, pd.DataFrame]:
        history_days = self.history_days
        div, price, train_size = div_price_train_size(tickers, end)
        div = div.iloc[train_size - history_days :]
        price = price.iloc[train_size - history_days :]

        return div, price


class ForecastParams(DataParams):
    """Метки не формируются, а признаки формируются только для последней даты."""

    def len(self, ticker) -> int:
        """Количество доступных примеров для данного тикера."""
        return 1

    def _div_price(self, tickers, end) -> Tuple[pd.DataFrame, pd.DataFrame]:
        history_days = self.history_days
        div, price, train_size = div_price_train_size(tickers, end)
        div = div.iloc[-history_days:]
        price = price.iloc[-history_days:]
        self._params["features"].pop("Label")

        return div, price
