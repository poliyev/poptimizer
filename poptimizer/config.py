"""Основные настраиваемые параметры"""
import logging
import pathlib

import pandas as pd


class POptimizerError(Exception):
    """Базовое исключение."""


# Конфигурация логгера
logging.basicConfig(level=logging.INFO)

# Количество колонок в распечатках без переноса на несколько страниц
pd.set_option("display.max_columns", 20)
pd.set_option("display.max_rows", 80)
pd.set_option("display.width", None)

# Путь к директории с данными
DATA_PATH = pathlib.Path(__file__).parents[1] / "data"

# Путь к директории с отчетам
REPORTS_PATH = pathlib.Path(__file__).parents[1] / "reports"

# Множитель, для переходя к после налоговым значениям
AFTER_TAX = 1 - 0.13

# Параметр для доверительных интервалов
T_SCORE = 2.0

# Максимальный объем одной торговой операции в долях портфеля
MAX_TRADE = 0.011

# Период в торговых днях, за который медианный оборот торгов
TURNOVER_PERIOD = 21 * 4

# Минимальный оборот - преимущества акции снижаются при приближении медианного оборота к данному уровню
TURNOVER_CUT_OFF = 3.6 * MAX_TRADE

# Параметры ML-модели
ML_PARAMS = {
    "data": (
        ("Label", {"days": 64, "on_off": True}),
        ("STD", {"days": 39, "on_off": True}),
        ("Ticker", {"on_off": True}),
        ("Mom12m", {"days": 242, "on_off": True, "periods": 1}),
        ("DivYield", {"days": 349, "on_off": True, "periods": 1}),
        ("Mom1m", {"days": 31, "on_off": False}),
        ("RetMax", {"days": 52, "on_off": True}),
        ("ChMom6m", {"days": 76, "on_off": True}),
    ),
    "model": {
        "bagging_temperature": 0.5350103284077592,
        "depth": 10,
        "l2_leaf_reg": 0.5944080741564565,
        "learning_rate": 0.003549670835757204,
        "one_hot_max_size": 100,
        "random_strength": 0.52921527733442,
    },
}
