Оптимизация долгосрочного портфеля акций
========================================
[![image](https://github.com/WLM1ke/poptimizer/workflows/tests/badge.svg)](https://github.com/WLM1ke/poptimizer/actions)
[![image](https://codecov.io/gh/WLM1ke/poptimizer/branch/master/graph/badge.svg)](https://codecov.io/gh/WLM1ke/poptimizer)

О проекте
---------

По образованию я человек далекий от программирования. Занимаюсь
инвестициями с 2008 года. Целью проекта является изучение
программирования и автоматизация процесса управления портфелем акций.

Используемый подход не предполагает баснословных доходностей, а нацелен
на получение результата чуть лучше рынка при рисках чуть меньше рынка
при относительно небольшом обороте. Портфель ценных бумаг должен быть
достаточно сбалансированным, чтобы его нестрашно было оставить без
наблюдения на продолжительное время.

Большинство частных инвесторов стремиться к быстрому обогащению и,
согласно известному афоризму Баффета, \"мало кто хочет разбогатеть
медленно\", поэтому проект является открытым. Стараюсь по возможности
исправлять ошибки, выявленные другими пользователями, и буду рад любой
помощи от более опытных программистов. Особенно приветствуются вопросы и
предложения по усовершенствованию содержательной части подхода к
управлению портфелем.

Проект находится в стадии развития и постоянно модифицируется (не всегда
удачно), поэтому может быть использован на свой страх и риск.

Основные особенности
--------------------

### Оптимизация портфеля

-   Базируется на [Modern portfolio
    theory](https://en.wikipedia.org/wiki/Modern_portfolio_theory)
-   При построении портфеля учитывается более 200 акций (включая
    иностранные) и ETF, обращающихся на MOEX
-   Используется ансамбль моделей для оценки неточности предсказаний
    ожидаемых доходностей и рисков отдельных активов
-   Используется робастная инкрементальная оптимизация на основе расчета
    достоверности улучшения метрик портфеля в результате торговли с
    учетом неточности имеющихся прогнозов вместо классической
    mean-variance оптимизации
-   Применяется [поправка
    Бонферрони](https://en.wikipedia.org/wiki/Bonferroni_correction) на
    множественное тестирование с учетом большое количества анализируемых
    активов

### Прогнозирование параметров активов

-   Используются нейронные сети на основе архитектуры
    [WaveNet](https://arxiv.org/abs/1609.03499) с большим receptive
    field для анализа длинных последовательностей котировок
-   Осуществляется совместное прогнозирование ожидаемой доходности и ее
    дисперсии с помощью подходов, базирующихся на [GluonTS:
    Probabilistic Time Series Models in
    Python](https://arxiv.org/abs/1906.05264)
-   Для моделирования толстых хвостов в распределениях доходностей
    применяются смеси логнормальных распределений
-   Используются устойчивые оценки исторических корреляционных матриц
    для большого числа активов с помощью сжатия
    [Ledoit-Wolf](http://www.ledoit.net/honey.pdf)

### Формирование ансамбля моделей

-   Осуществляется выбор моделей из многомерного пространства
    гиперпараметров сетей, их оптимизаторов и комбинаций признаков
-   Для исследования пространства применяются подходы алгоритма
    [Имитации
    отжига](https://en.wikipedia.org/wiki/Simulated_annealing)
-   Для масштабирования локальной области поиска и кодирования
    гиперпараметров используются принципы [дифференциальной
    эволюции](https://en.wikipedia.org/wiki/Differential_evolution)
-   Для выбора моделей в локальной области применяется распределение
    [Коши](https://en.wikipedia.org/wiki/Cauchy_distribution) для
    осуществления редких не локальных прыжков в пространстве
    гиперпараметров
-   При отборе претендентов в ансамбль осуществляется [последовательное
    тестирование](https://en.wikipedia.org/wiki/Sequential_analysis#Alpha_spending_functions)
    с соответствующими корректировками [уровней
    значимости](https://arxiv.org/abs/1906.09712)

### Источники данных

-   Реализована загрузка котировок всех акций (включая иностранные) и
    ETF, обращающихся на MOEX
-   Поддерживается в актуальном состоянии база данных дивидендов с 2015г
    по включенным в анализ акциям
-   Реализована возможность сверки базы данных дивидендов с информацией
    на сайтах:

> -   [www.dohod.ru](https://www.dohod.ru/ik/analytics/dividend)
> -   [www.conomy.ru](https://www.conomy.ru/dates-close/dates-close2)
> -   [bcs-express.ru](https://bcs-express.ru/dividednyj-kalendar)
> -   [www.smart-lab.ru](https://smart-lab.ru/dividends/index/order_by_yield/desc/)
> -   [закрытияреестров.рф](https://закрытияреестров.рф/)
> -   [finrange.com](https://finrange.com/)
> -   [investmint.ru](https://investmint.ru/)
> -   [www.nasdaq.com](https://www.nasdaq.com/)
> -   [www.streetinsider.com](https://www.streetinsider.com/)

Направления дальнейшего развития
--------------------------------

- Применение нелинейного сжатия Ledoit-Wolf для оценки корреляции
    активов
- Реализация сервиса на Go для загрузки всей необходимой информации
- Рефакторинг кода на основе
    [DDD](https://en.wikipedia.org/wiki/Domain-driven_design),
    [MyPy](http://mypy.readthedocs.org/en/latest/) и
    [wemake](https://wemake-python-stylegui.de/en/latest/)
- Использование архитектур на основе
    [трансформеров](https://en.wikipedia.org/wiki/Transformer_(machine_learning_model))
    вместо WaveNet
- Поиск оптимальной архитектуры сетей с помощью эволюции с \"нуля\" по
    аналогии с [Evolving Neural Networks through Augmenting
    Topologies](http://nn.cs.utexas.edu/downloads/papers/stanley.ec02.pdf)
- Использование Reinforcement learning для построения портфеля

FAQ
---

### Какие инструменты нужны для запуска программы?

Последняя версия Python, все зависимости из [requirements.txt](https://github.com/WLM1ke/poptimizer/blob/18756e8bdbfcac3ebd7ba241f86b25bcb27cc22f/requirements.txt),
MongoDB и MongoDB Database Tools.

### Как запускать программу?

Запуск реализован через CLI:

`python3 -m poptimizer`

После этого можно посмотреть перечень команд и help к ним, а дальше самому разбираться в коде.
Основные команды для запуска описаны в файле [\_\_main\_\_.py](https://github.com/WLM1ke/poptimizer/blob/18756e8bdbfcac3ebd7ba241f86b25bcb27cc22f/poptimizer/__main__.py).
Сначала необходимо запустить функцию `evolve` для обучения моделей. После этого можно запустить `optimize` для 
оптимизации портфеля.

### Есть ли у программы какие-нибудь настройки?

Настройки описаны в файле [config.template](https://github.com/WLM1ke/poptimizer/blob/18756e8bdbfcac3ebd7ba241f86b25bcb27cc22f/config/config.template).
При отсутствии файла конфигурации будут использоваться значения по умолчанию. 

### Как ввести свой портфель?

Пример заполнения файла с портфеля с базовым набором бумаг содержится в файле [base.yaml](https://github.com/WLM1ke/poptimizer/blob/18756e8bdbfcac3ebd7ba241f86b25bcb27cc22f/portfolio/base.yaml).
В этом каталоге можно хранить множество файлов (например по отдельным брокерским счетам), информация из них будет 
объединяться в единый портфель. 

### У меня в портфеле не так много бумаг - нулевые значения по множеству позиций как-нибудь влияют на эволюцию?

Для эволюции принципиален перечень бумаг, а не их количество. 

### У меня в портфеле не так много бумаг - можно оставить только свои?

При желании можно сократить, количество бумаг, но их должно быть не меньше половины из базового набора, чтобы получалось 
достаточно большое количество обучающих примеров для тренировки моделей.

### В моем портфеле есть бумаги, отсутствующие в базовом наборе - можно ли их добавить?

Можно добавить любые акции, включая иностранные, и ETF, обращающиеся на MOEX. Для корректной работы так же может 
потребоваться дополнить базу по дивидендам, если они выплачивались с 2015 года.
 
### Что отражают LOWER и UPPER в разделе оптимизация портфеля?

Нижняя и верхняя граница доверительного интервала влияния указанной бумаги на качество портфеля. Если верхняя граница 
одной бумаги ниже нижней границы второй бумаги, то целесообразно сокращать позицию по первой бумаге и наращивать позицию 
по второй. При выдаче рекомендаций дополнительно учитывается, что зазор между границами должен покрывать транзакционные 
издержки.

Особые благодарности
--------------------

-   [Evgeny Pogrebnyak](https://github.com/epogrebnyak) за помощь в
    освоении Python
-   [RomaKoks](https://github.com/RomaKoks) за полезные советы по
    автоматизации некоторых этапов работы программы и исправлению ошибок
-   [AlexQww](https://github.com/AlexQww) за содержательные обсуждения
    подходов к управлению портфелем, которые стали катализатором
    множества изменений в программе
