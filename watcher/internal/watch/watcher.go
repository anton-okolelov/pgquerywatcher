package watch

import (
	"context"
	"strings"
	"time"

	"github.com/anton.okolelov/pgquerywatcher/internal/config"
	"github.com/anton.okolelov/pgquerywatcher/internal/pg"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog"
)

type Watcher struct {
	targetDb *pgxpool.Pool
	log      zerolog.Logger
}

type PgStatStatementsRow struct {
	Query     string  `db:"query"`
	Calls     uint64  `db:"calls"`
	TotalTime float64 `db:"total_exec_time"`
}

type QueryStats struct {
	Time           time.Time `db:"time"`
	Query          string    `db:"query"`
	ExecMeanTimeMs float64   `db:"exec_mean_time_ms"`
}

func NewWatcher(targetDbConfig config.DbConfig, log zerolog.Logger) (*Watcher, error) {
	targetDb, err := pg.NewDB(targetDbConfig)
	if err != nil {
		return nil, err
	}

	return &Watcher{targetDb: targetDb, log: log}, err
}

func (w Watcher) Watch() error {

	previousStatStatements, err := w.getStatStatements()
	if err != nil {
		return err
	}

	for true {
		time.Sleep(3 * time.Second)
		currentStatStatements, err := w.getStatStatements()

		if err != nil {
			return err
		}

		// записываем изменения
		w.logChanges(currentStatStatements, previousStatStatements)
		previousStatStatements = currentStatStatements
	}
	return nil
}

func (w Watcher) logChanges(currentStatStatements map[string]PgStatStatementsRow, previousStatStatements map[string]PgStatStatementsRow) {
	currTime := time.Now()
	for query, curr := range currentStatStatements {
		if strings.Contains(query, "pg_stat_statements") {
			continue
		}
		prev, exists := previousStatStatements[query]
		if !exists || prev.Calls > curr.Calls {
			prev = PgStatStatementsRow{Calls: 0, TotalTime: 0}
		}

		if curr.Calls == prev.Calls {
			continue
		}

		timeDiff := curr.TotalTime - prev.TotalTime
		countDiff := curr.Calls - prev.Calls
		execMeanTime := timeDiff / float64(countDiff)

		w.log.Info().
			Bool("is_sql_stats", true).
			Time("time", currTime).
			Float64("time_diff", timeDiff).
			Uint64("count_diff", countDiff).
			Float64("exec_mean_time", execMeanTime).
			Msg(query)
	}
}

func (w Watcher) getStatStatements() (map[string]PgStatStatementsRow, error) {

	rows, err := w.targetDb.Query(context.Background(), `
		SELECT query, calls, total_exec_time 
			FROM pg_stat_statements`,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	queriesStats := map[string]PgStatStatementsRow{}
	for rows.Next() {
		var queryStat PgStatStatementsRow
		err = rows.Scan(&queryStat.Query, &queryStat.Calls, &queryStat.TotalTime)
		if err != nil {
			return nil, err
		}
		queriesStats[queryStat.Query] = queryStat
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return queriesStats, nil
}
