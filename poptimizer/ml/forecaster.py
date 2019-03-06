"""Составляет ML-прогноз."""
from typing import Tuple

import catboost
import numpy as np
import pandas as pd

from poptimizer import data
from poptimizer.config import POptimizerError, ML_PARAMS
from poptimizer.ml import examples, ledoit_wolf, cv
from poptimizer.ml.feature import YEAR_IN_TRADING_DAYS
from poptimizer.portfolio.metrics import Forecast


def cv_results(cases: examples.Examples, params):
    """Получает необходимые результаты кросс-валидации."""
    result = cv.valid_model(params, cases)
    return result["std"], result["r2"], result["params"]


def fit_clf(valid_params: tuple, cases: examples.Examples):
    """Тренирует ML-модель на основе параметров с учетом количества итераций."""
    data_params, model_params = valid_params
    train_params, predict_params = cases.train_predict_pool_params()
    learn_pool = catboost.Pool(**train_params)
    clf = catboost.CatBoostRegressor(**model_params)
    clf.fit(learn_pool)
    return clf, learn_pool.num_row(), predict_params


def predict_mean(clf, predict_pool_params):
    """Прогноз ожидаемой доходности."""
    predict_pool = catboost.Pool(**predict_pool_params)
    raw_prediction = clf.predict(predict_pool)
    scaler = predict_pool_params["data"].iloc[:, 0]
    return raw_prediction * scaler.values


def validate_cov(cov, predict_pool_params):
    """Проверяет совпадение ковариации с использовавшейся для нормирования."""
    scaler = predict_pool_params["data"].iloc[:, 0]
    if not np.allclose(np.diag(cov), scaler.values ** 2):
        raise POptimizerError(
            "Расчетная ковариация не совпадает с использовавшейся для нормирования"
        )


def ledoit_wolf_cov(predict_pool_params, cv_params, tickers, date, ml_std):
    """Ковариационная матрица на основе Ledoit Wolf и вспомогательные данные.

    Оригинальная матрица корректируется в сторону не смещенной оценки на малой выборке и точность
    ML-прогноза.
    """
    mean_days, scaler_days = cv_params[0][0][1]["days"], cv_params[0][1][1]["days"]
    returns = data.log_total_returns(tickers, date)
    returns = returns.iloc[-scaler_days:,]
    cov, average_cor, shrinkage = ledoit_wolf.shrinkage(returns.values)
    cov *= scaler_days / (scaler_days - 1)
    validate_cov(cov, predict_pool_params)
    cov *= ml_std ** 2 * mean_days
    return cov, average_cor, shrinkage


# noinspection PyUnresolvedReferences
def make_forecast(
    tickers: Tuple[str, ...], date: pd.Timestamp, params=ML_PARAMS
) -> Forecast:
    """Создает прогноз для набора тикеров на указанную дату.

    :param tickers:
        Тикеры, для которых необходимо составить прогноз.
    :param date:
        Дата, на которую нужен прогноз.
    :param params:
        Параметры ML-модели для прогноза.
    :return:
        Прогнозная доходность, ковариация и дополнительная информация.
    """
    cases = examples.Examples(tickers, date, params[0])
    ml_std, r2, cv_params = cv_results(cases, params)
    clf, num_cases, predict_params = fit_clf(cv_params, cases)
    feature_importance = pd.Series(
        clf.feature_importances_, cases.get_features_names(), name="Importance"
    )
    mean = predict_mean(clf, predict_params)
    cov, average_cor, shrinkage = ledoit_wolf_cov(
        predict_params, cv_params, tickers, date, ml_std
    )
    return Forecast(
        date=date,
        tickers=tickers,
        mean=mean * YEAR_IN_TRADING_DAYS,
        cov=cov * YEAR_IN_TRADING_DAYS,
        num_cases=num_cases,
        trees=cv_params[1]["iterations"],
        depth=cv_params[1]["depth"],
        feature_importance=feature_importance,
        r2=r2,
        average_cor=average_cor,
        shrinkage=shrinkage,
    )
