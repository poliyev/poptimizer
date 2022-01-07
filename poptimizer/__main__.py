"""Запуск основных операций с помощью CLI."""
import logging

import typer
import os

from poptimizer import config
from poptimizer.data.views import div_status
from poptimizer.evolve import Evolution
from poptimizer.portfolio import load_from_yaml, optimizer_hmean, optimizer_resample

LOGGER = logging.getLogger()

os.environ["KMP_DUPLICATE_LIB_OK"] = "TRUE"


def evolve() -> None:
    """Run evolution."""
    ev = Evolution()
    ev.evolve()


def dividends(ticker: str) -> None:
    """Get dividends status."""
    div_status.dividends_validation(ticker)


def optimize(date: str = typer.Argument(..., help="YYYY-MM-DD")) -> None:
    """Optimize portfolio."""
    port = load_from_yaml(date)
    opt_type = {
        "resample": optimizer_resample.Optimizer,
        "hmean": optimizer_hmean.Optimizer,
    }[config.OPTIMIZER]
    opt = opt_type(port)
    LOGGER.info(opt.portfolio)
    LOGGER.info(opt.metrics)
    LOGGER.info(opt)
    div_status.new_dividends(tuple(port.index[:-2]))


if __name__ == "__main__":
    app = typer.Typer(help="Run poptimizer subcommands.", add_completion=False)

    app.command()(evolve)
    app.command()(dividends)
    app.command()(optimize)

    app(prog_name="poptimizer")
