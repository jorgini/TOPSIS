package repository

import (
	"context"
	"errors"
	"fmt"
	"time"
	"webApp/configs"
	"webApp/entity"
)

type TaskDao struct {
	c   IConnectionFactory
	cfg *configs.DbConfig
}

func NewTaskDao(factory IConnectionFactory, config *configs.DbConfig) *TaskDao {
	return &TaskDao{
		c:   factory,
		cfg: config,
	}
}

func (t *TaskDao) CreateNewTask(ctx context.Context, task *entity.TaskModel) (int64, error) {
	query := fmt.Sprintf(`INSERT INTO %s (maintainer, password, title, description, last_change, 
		task_type, method, calc_settings, ling_scale, status) values
		($1, '', $2, $3, $4, $5, $6, $7, $8, $9) RETURNING sid`, t.cfg.TaskTable)

	conn := t.c.getConnection()
	if conn == nil {
		return 0, errors.New("cant connect to db")
	}

	var sid int64
	row := conn.QueryRowxContext(ctx, query, task.MaintainerID, task.Title, task.Description, time.Now(),
		task.TaskType, task.Method, task.CalcSettings, task.LingScale, task.Status)
	if err := row.Scan(&sid); err != nil {
		return 0, errors.Join(err, t.c.closeConnection())
	}
	return sid, t.c.closeConnection()
}

func (t *TaskDao) ValidateUser(ctx context.Context, uid, sid int64) error {
	query := fmt.Sprintf("SELECT maintainer FROM %s WHERE sid=$1", t.cfg.TaskTable)

	conn := t.c.getConnection()
	if conn == nil {
		return errors.New("cant connect to db")
	}

	var trueUid int64
	row := conn.QueryRowxContext(ctx, query, sid)
	if err := row.Scan(&trueUid); err != nil {
		return errors.Join(err, t.c.closeConnection())
	}
	if trueUid != uid {
		return errors.Join(errors.New("invalid operation for current user"), t.c.closeConnection())
	}
	return t.c.closeConnection()
}

func (t *TaskDao) UpdateTask(ctx context.Context, sid int64, input *entity.TaskModel) error {
	query := fmt.Sprintf(`UPDATE %s SET title=$1, description=$2, last_change=$3, task_type=$4,
		method=$5, calc_settings=$6, ling_scale=$7, status=$8 WHERE sid=$9`, t.cfg.TaskTable)

	conn := t.c.getConnection()
	if conn == nil {
		return errors.New("cant connect to db")
	}

	_, err := conn.ExecContext(ctx, query, input.Title, input.Description, time.Now(), input.TaskType,
		input.Method, input.CalcSettings, input.LingScale, entity.Draft, sid)
	if err != nil {
		return errors.Join(err, t.c.closeConnection())
	}

	return t.c.closeConnection()
}

func (t *TaskDao) GetTask(ctx context.Context, sid int64) (*entity.TaskModel, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE sid=$1", t.cfg.TaskTable)

	conn := t.c.getConnection()
	if conn == nil {
		return nil, errors.New("cant connect to db")
	}

	var task entity.TaskModel
	if err := conn.GetContext(ctx, &task, query, sid); err != nil {
		return nil, errors.Join(err, t.c.closeConnection())
	}
	return &task, t.c.closeConnection()
}

func (t *TaskDao) DeleteTask(ctx context.Context, sid int64) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE sid=$1", t.cfg.TaskTable)

	conn := t.c.getConnection()
	if conn == nil {
		return errors.New("cant connect to db")
	}

	if _, err := conn.ExecContext(ctx, query, sid); err != nil {
		return errors.Join(err, t.c.closeConnection())
	}
	return t.c.closeConnection()
}

func (t *TaskDao) SetPassword(ctx context.Context, sid int64, password string) error {
	query := fmt.Sprintf("UPDATE %s SET password=$1 WHERE sid=$2", t.cfg.TaskTable)

	conn := t.c.getConnection()
	if conn == nil {
		return errors.New("cant connect to db")
	}

	if _, err := conn.ExecContext(ctx, query, password, sid); err != nil {
		return errors.Join(err, t.c.closeConnection())
	}
	return t.c.closeConnection()
}

func (t *TaskDao) SetCriteria(ctx context.Context, sid int64, criteria entity.Criteria) error {
	query := fmt.Sprintf("UPDATE %s SET criteria=$1, last_change=$2 WHERE sid=$3", t.cfg.TaskTable)

	conn := t.c.getConnection()
	if conn == nil {
		return errors.New("cant connect to db")
	}

	if _, err := conn.ExecContext(ctx, query, criteria, time.Now(), sid); err != nil {
		return errors.Join(err, t.c.closeConnection())
	}
	return t.c.closeConnection()
}

func (t *TaskDao) SetAlts(ctx context.Context, sid int64, alts entity.Alts) error {
	query := fmt.Sprintf("UPDATE %s SET alternatives=$1, last_change=$2 WHERE sid=$3", t.cfg.TaskTable)

	conn := t.c.getConnection()
	if conn == nil {
		return errors.New("cant connect to db")
	}

	if _, err := conn.ExecContext(ctx, query, alts, time.Now(), sid); err != nil {
		return errors.Join(err, t.c.closeConnection())
	}
	return t.c.closeConnection()
}

func (t *TaskDao) GetCriteria(ctx context.Context, sid int64) (entity.Criteria, error) {
	query := fmt.Sprintf("SELECT criteria FROM %s WHERE sid=$1", t.cfg.TaskTable)

	conn := t.c.getConnection()
	if conn == nil {
		return nil, errors.New("cant connect to db")
	}

	var criteria entity.Criteria
	if err := conn.GetContext(ctx, &criteria, query, sid); err != nil {
		return nil, errors.Join(err, t.c.closeConnection())
	}
	return criteria, t.c.closeConnection()
}

func (t *TaskDao) GetAlts(ctx context.Context, sid int64) (entity.Alts, error) {
	query := fmt.Sprintf("SELECT alternatives FROM %s WHERE sid=$1", t.cfg.TaskTable)

	conn := t.c.getConnection()
	if conn == nil {
		return nil, errors.New("cant connect to db")
	}

	var alts entity.Alts
	if err := conn.GetContext(ctx, &alts, query, sid); err != nil {
		return nil, errors.Join(err, t.c.closeConnection())
	}
	return alts, t.c.closeConnection()
}

func (t *TaskDao) UpdateCriteria(ctx context.Context, sid int64, criteria entity.Criteria) error {
	query := fmt.Sprintf("UPDATE %s SET criteria=$1, last_change=$2 WHERE sid=$3", t.cfg.TaskTable)

	conn := t.c.getConnection()
	if conn == nil {
		return errors.New("cant connect to db")
	}

	if _, err := conn.ExecContext(ctx, query, criteria, time.Now(), sid); err != nil {
		return errors.Join(err, t.c.closeConnection())
	}

	return t.c.closeConnection()
}

func (t *TaskDao) UpdateAlts(ctx context.Context, sid int64, alts entity.Alts) error {
	query := fmt.Sprintf("UPDATE %s SET alternatives=$1, last_change=$2 WHERE sid=$3", t.cfg.TaskTable)

	conn := t.c.getConnection()
	if conn == nil {
		return errors.New("cant connect to db")
	}

	if _, err := conn.ExecContext(ctx, query, alts, time.Now(), sid); err != nil {
		return errors.Join(err, t.c.closeConnection())
	}
	return t.c.closeConnection()
}

func (t *TaskDao) GetAllSolutions(ctx context.Context, uid int64) ([]entity.TaskModel, error) {
	query := fmt.Sprintf("SELECT t.* FROM %s t LEFT JOIN %s m ON t.sid=m.sid WHERE m.uid=$1 OR maintainer=$1",
		t.cfg.TaskTable, t.cfg.MatrixTable)

	conn := t.c.getConnection()
	if conn == nil {
		return nil, errors.New("cant connect to db")
	}

	var tasks []entity.TaskModel
	if err := conn.SelectContext(ctx, &tasks, query, uid); err != nil {
		return nil, errors.Join(err, t.c.closeConnection())
	}
	return tasks, t.c.closeConnection()
}

func (t *TaskDao) ConnectToTask(ctx context.Context, sid int64, password string) error {
	query := fmt.Sprintf("SELECT password FROM %s WHERE sid=$1", t.cfg.TaskTable)

	conn := t.c.getConnection()
	if conn == nil {
		return errors.New("cant connect to db")
	}

	var truePass string
	row := conn.QueryRowxContext(ctx, query, sid)
	if err := row.Scan(&truePass); err != nil {
		return errors.Join(err, t.c.closeConnection())
	}
	if truePass != password {
		return errors.Join(errors.New("incorrect password"), t.c.closeConnection())
	}
	return t.c.closeConnection()
}

func (t *TaskDao) SetLastChange(ctx context.Context, sid int64) error {
	query := fmt.Sprintf("UPDATE %s SET last_change=$1 WHERE sid=$2", t.cfg.TaskTable)

	conn := t.c.getConnection()
	if conn == nil {
		return errors.New("cant connect to db")
	}

	if _, err := conn.ExecContext(ctx, query, time.Now(), sid); err != nil {
		return errors.Join(err, t.c.closeConnection())
	}
	return t.c.closeConnection()
}

func (t *TaskDao) SetExpertsWeights(ctx context.Context, sid int64, weights entity.Weights) error {
	query := fmt.Sprintf("UPDATE %s SET experts_weights=$1, last_change=$2 WHERE sid=$3", t.cfg.TaskTable)

	conn := t.c.getConnection()
	if conn == nil {
		return errors.New("cant connect to db")
	}

	if _, err := conn.ExecContext(ctx, query, weights, time.Now(), sid); err != nil {
		return errors.Join(err, t.c.closeConnection())
	}
	return t.c.closeConnection()
}
