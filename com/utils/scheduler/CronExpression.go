/*
@Time : 2018/12/21 13:33 
@Author : sky
@Description:
@File : CronExpression
@Software: GoLand
*/
package scheduler

const(

	 SECOND =  iota
	 MINUTE
	 HOUR
	 DAY_OF_MONTH
	 MONTH
	 DAY_OF_WEEK
	 YEAR
	 ALL_SPEC_INT=99
	 NO_SPEC_INT=98
)

var monthMap=map[string]int{
	"JAN":0,
	"FEB":1,
	"MAR":2,
	"APR":3,
	"MAY":4,
	"JUN":5,
	"JUL":6,
	"AUG":7,
	"SEP":8,
	"OCT":9,
	"NOV":10,
	"DEC":11}
var dayMap =map[string]int{
	"SUN":1,
	"MON":2,
	"TUE":3,
	"WED":4,
	"THU":5,
	"FRI":6,
	"SAT":7}

type CronExpression struct {

}

func NewCronExpression()(CronExpression,error){



	return CronExpression{},nil
}

