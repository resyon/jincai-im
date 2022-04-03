package conf

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetMysqlDSN(t *testing.T) {
	dsn := GetMysqlDSN()
	assert.NotEqual(t, mysqlDsnFmt, dsn, "fail to parse mysql conf or conf is empty")
}
