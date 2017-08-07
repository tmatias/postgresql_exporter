package gauges

import "github.com/prometheus/client_golang/prometheus"

func (g *Gauges) Up() prometheus.Gauge {
	lbl := prometheus.Labels{}
	for k, v := range g.labels {
		lbl[k] = v
	}
	lbl["version"] = g.version()
	return g.new(
		prometheus.GaugeOpts{
			Name:        "postgresql_up",
			Help:        "Dabatase is up and accepting connections",
			ConstLabels: lbl,
		},
		"SELECT 1",
	)
}

func (g *Gauges) Size() prometheus.Gauge {
	return g.new(
		prometheus.GaugeOpts{
			Name:        "postgresql_size_bytes",
			Help:        "Dabatase size in bytes",
			ConstLabels: g.labels,
		},
		"SELECT pg_database_size(current_database())",
	)
}

func (g *Gauges) TempSize() prometheus.Gauge {
	return g.new(
		prometheus.GaugeOpts{
			Name:        "postgresql_temp_bytes",
			Help:        "Temp size in bytes",
			ConstLabels: g.labels,
		},
		"SELECT temp_bytes FROM pg_stat_database WHERE datname = current_database()",
	)
}

func (g *Gauges) TempFiles() prometheus.Gauge {
	return g.new(
		prometheus.GaugeOpts{
			Name:        "postgresql_temp_files",
			Help:        "Count of temp files",
			ConstLabels: g.labels,
		},
		"SELECT temp_files FROM pg_stat_database WHERE datname = current_database()",
	)
}

func (g *Gauges) Deadlocks() prometheus.Gauge {
	return g.new(
		prometheus.GaugeOpts{
			Name:        "postgresql_deadlocks",
			Help:        "Number of deadlocks in the last 2m",
			ConstLabels: g.labels,
		},
		`
			SELECT count(*) FROM pg_locks bl
			JOIN pg_stat_activity a
			ON a.pid = bl.pid JOIN pg_locks kl
			ON kl.transactionid = bl.transactionid
			AND kl.pid != bl.pid JOIN pg_stat_activity ka
			ON ka.pid = kl.pid WHERE NOT bl.granted
			AND (ka.query_start >= (now() - interval '2 minutes')
			OR a.query_start >= (now() - interval '2 minutes'))
			AND a.datname = current_database()
		`,
	)
}
