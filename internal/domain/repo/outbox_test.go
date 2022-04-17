package repo

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/ssup2ket/ssup2ket-auth-service/internal/domain/entity"
	"github.com/ssup2ket/ssup2ket-auth-service/internal/test"
)

func TestOutbox(t *testing.T) {
	suite.Run(t, new(outboxSuite))
}

type outboxSuite struct {
	suite.Suite
	sqlMock sqlmock.Sqlmock

	tx   *DBTxImp
	repo OutboxRepo
}

func (o *outboxSuite) SetupTest() {
	var err error
	var db *sql.DB

	// Init sqlMock
	db, o.sqlMock, err = sqlmock.New()
	require.NoError(o.T(), err)

	// Init DB
	primaryMySQL, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}))
	require.NoError(o.T(), err)

	// Init transaction, repo
	o.tx = NewDBTxImp(primaryMySQL)
	o.repo = NewOutboxRepoImp(primaryMySQL)
}

func (o *outboxSuite) AfterTest(_, _ string) {
	require.NoError(o.T(), o.sqlMock.ExpectationsWereMet())
}

func (o *outboxSuite) TestCreateSuccess() {
	o.sqlMock.ExpectBegin()
	o.sqlMock.ExpectExec(regexp.QuoteMeta("INSERT INTO `outboxes` (`id`,`created_at`,`aggregatetype`,`aggregateid`,`eventtype`,`payload`,`spancontext`) VALUES (?,?,?,?,?,?,?)")).
		WithArgs(test.OutboxIDCorrect, sqlmock.AnyArg(), test.OutboxAggregateTypeCorrect, test.OutboxAggregateIDCorrect, test.OutboxEventTypeCorrect, test.OutboxPayloadCorrect, test.OutboxSpanContextCorrect).
		WillReturnResult(sqlmock.NewResult(1, 1))
	o.sqlMock.ExpectCommit()

	err := o.repo.Create(context.Background(), &entity.Outbox{
		ID:            test.OutboxIDCorrect,
		AggregateType: test.OutboxAggregateTypeCorrect,
		AggregateID:   test.OutboxAggregateIDCorrect,
		EventType:     test.OutboxEventTypeCorrect,
		Payload:       test.OutboxPayloadCorrect,
		SpanContext:   test.OutboxSpanContextCorrect,
	})
	require.NoError(o.T(), err)
}

func (o *outboxSuite) TestCreateError() {
	o.sqlMock.ExpectBegin()
	o.sqlMock.ExpectExec(regexp.QuoteMeta("INSERT INTO `outboxes` (`id`,`created_at`,`aggregatetype`,`aggregateid`,`eventtype`,`payload`,`spancontext`) VALUES (?,?,?,?,?,?,?)")).
		WithArgs(test.OutboxIDCorrect, sqlmock.AnyArg(), test.OutboxAggregateTypeCorrect, test.OutboxAggregateIDCorrect, test.OutboxEventTypeCorrect, test.OutboxPayloadCorrect, test.OutboxSpanContextCorrect).
		WillReturnError(fmt.Errorf("error"))
	o.sqlMock.ExpectRollback()

	err := o.repo.Create(context.Background(), &entity.Outbox{
		ID:            test.OutboxIDCorrect,
		AggregateType: test.OutboxAggregateTypeCorrect,
		AggregateID:   test.OutboxAggregateIDCorrect,
		EventType:     test.OutboxEventTypeCorrect,
		Payload:       test.OutboxPayloadCorrect,
		SpanContext:   test.OutboxSpanContextCorrect,
	})
	require.Error(o.T(), err)
}

func (o *outboxSuite) TestCreateWithTxSuccess() {
	o.sqlMock.ExpectBegin()
	o.sqlMock.ExpectExec(regexp.QuoteMeta("INSERT INTO `outboxes` (`id`,`created_at`,`aggregatetype`,`aggregateid`,`eventtype`,`payload`,`spancontext`) VALUES (?,?,?,?,?,?,?)")).
		WithArgs(test.OutboxIDCorrect, sqlmock.AnyArg(), test.OutboxAggregateTypeCorrect, test.OutboxAggregateIDCorrect, test.OutboxEventTypeCorrect, test.OutboxPayloadCorrect, test.OutboxSpanContextCorrect).
		WillReturnResult(sqlmock.NewResult(1, 1))
	o.sqlMock.ExpectCommit()

	tx, _ := o.tx.Begin()
	err := o.repo.WithTx(tx).Create(context.Background(), &entity.Outbox{
		ID:            test.OutboxIDCorrect,
		AggregateType: test.OutboxAggregateTypeCorrect,
		AggregateID:   test.OutboxAggregateIDCorrect,
		EventType:     test.OutboxEventTypeCorrect,
		Payload:       test.OutboxPayloadCorrect,
		SpanContext:   test.OutboxSpanContextCorrect,
	})
	require.NoError(o.T(), err)
	tx.Commit()
}

func (o *outboxSuite) TestDeleteSuccess() {
	o.sqlMock.ExpectBegin()
	o.sqlMock.ExpectExec(regexp.QuoteMeta("DELETE FROM `outboxes` WHERE id = ?")).
		WithArgs(test.OutboxIDCorrect).
		WillReturnResult(sqlmock.NewResult(1, 1))
	o.sqlMock.ExpectCommit()

	err := o.repo.Delete(context.Background(), test.OutboxIDCorrect)
	require.NoError(o.T(), err)
}

func (o *outboxSuite) TestDeleteError() {
	o.sqlMock.ExpectBegin()
	o.sqlMock.ExpectExec(regexp.QuoteMeta("DELETE FROM `outboxes` WHERE id = ?")).
		WithArgs(test.OutboxIDCorrect).
		WillReturnError(fmt.Errorf("error"))
	o.sqlMock.ExpectRollback()

	err := o.repo.Delete(context.Background(), test.OutboxIDCorrect)
	require.Error(o.T(), err)
}
