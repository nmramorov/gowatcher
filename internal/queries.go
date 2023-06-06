package metrics

// Запросы, используемые при работе приложения.
const (
	CreateGaugeTable string = `CREATE TABLE IF NOT EXISTS gaugeMetrics (
		_id TEXT,
		mtype TEXT,
		_value DOUBLE PRECISION,
	  	date TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`
	CreateCounterTable string = `CREATE TABLE IF NOT EXISTS counterMetrics (
		_id TEXT,
		mtype TEXT,
		_value INTEGER,
		date TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`
	InsertIntoGauge          = `INSERT INTO gaugemetrics VALUES ($1, $2, $3);`
	InsertIntoCounter        = `INSERT INTO countermetrics VALUES ($1, $2, $3);`
	SelectFromGauge          = `SELECT * FROM gaugemetrics WHERE _id=$1 ORDER BY date DESC LIMIT 1`
	SelectFromCounter string = `SELECT * FROM countermetrics WHERE _id=$1 ORDER BY date DESC LIMIT 1`
)
