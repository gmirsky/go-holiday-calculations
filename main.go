package main

import (
	"crypto/rand"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"gopkg.in/mgo.v2/bson"
	"gopkg.in/yaml.v2"

	"github.com/soniakeys/meeus/julian"
	"github.com/soniakeys/meeus/solstice"
)

const moveToMonday bool = true
const dontMoveToMonday bool = false
const outputFilePermission os.FileMode = os.FileMode(0644) // ---RW-R--R--
var waitGroup sync.WaitGroup

// Note that structure type names need to be capitiazized to be exported using the
// marshall functions

// NLHolidays  Netherlands holiday structure
type NLHolidays struct {
	Nieuwjaardag    time.Time `json:"Nieuwjaardag" yaml:"Nieuwjaardag" bson:"Nieuwjaardag"`          //Nieuwjaardag    = NewYear
	Goedevrijdag    time.Time `json:"Goedevrijdag" yaml:"Goedevrijdag" bson:"Goedevrijdag"`          //Goedevrijdag    = GoodFriday
	Paasmaandag     time.Time `json:"Paasmaandag" yaml:"Paasmaandag" bson:"Paasmaandag"`             //Paasmaandag     = EasterMonday
	Koningsdag      time.Time `json:"Koningsdag" yaml:"Koningsdag" bson:"Koningsdag"`                //Koningsdag      = Kings day
	Bevrijdingsdag  time.Time `json:"Bevrijdingsdag" yaml:"Bevrijdingsdag" bson:"Bevrijdingsdag"`    //Bevrijdingsdag  = May, 5
	Hemelvaart      time.Time `json:"Hemelvaart" yaml:"Hemelvaart" bson:"Hemelvaart"`                //Hemelvaart      = DE_Himmelfahrt
	Pinkstermaandag time.Time `json:"Pinkstermaandag" yaml:"Pinkstermaandag" bson:"Pinkstermaandag"` //Pinkstermaandag = DE_Pfingstmontag
	Eerstekerstdag  time.Time `json:"Eerstekerstdag" yaml:"Eerstekerstdag" bson:"Eerstekerstdag"`    //Eerstekerstdag  = Christmas
	Tweedekerstdag  time.Time `json:"Tweedekerstdag" yaml:"Tweedekerstdag" bson:"Tweedekerstdag"`    //Tweedekerstdag  = Christmas2
}

// UKHolidays UK holiday sturcture
type UKHolidays struct {
	NewYearsDay   time.Time `json:"NewYearsDay" yaml:"NewYearsDay" bson:"NewYearsDay"`
	GoodFriday    time.Time `json:"GoodFriday" yaml:"GoodFriday" bson:"GoodFriday"`
	EasterMonday  time.Time `json:"EasterMonday" yaml:"EasterMonday" bson:"EasterMonday"`
	EarlyMay      time.Time `json:"EarlyMay" yaml:"EarlyMay" bson:"EarlyMay"`
	SpringHoliday time.Time `json:"SpringHoliday" yaml:"SpringHoliday" bson:"SpringHoliday"`
	ChristmasDay  time.Time `json:"ChristmasDay" yaml:"ChristmasDay" bson:"ChristmasDay"`
	BoxingDay     time.Time `json:"BoxingDay" yaml:"BoxingDay" bson:"BoxingDay"`
}

// ECBTarget2Holidays ECB Target 2 holiday structure
type ECBTarget2Holidays struct {
	NewYearsDay      time.Time `json:"NewYearsDay" yaml:"NewYearsDay" bson:"NewYearsDay"`
	GoodFriday       time.Time `json:"GoodFriday" yaml:"GoodFriday" bson:"GoodFriday"`
	EasterMonday     time.Time `json:"EasterMonday" yaml:"EasterMonday" bson:"EasterMonday"`
	LaborDay         time.Time `json:"LaborDay" yaml:"LaborDay" bson:"LaborDay"`
	ChristmasDay     time.Time `json:"ChristmasDay" yaml:"ChristmasDay" bson:"ChristmasDay"`
	ChristmasHoliday time.Time `json:"ChristmasHoliday" yaml:"ChristmasHoliday" bson:"ChristmasHoliday"`
}

// EU holiday structure
// type EUHolidays struct {
// 	NewYearsDay       time.Time `json:"NewYearsDay" yaml:"NewYearsDay" bson:"NewYearsDay"`
// 	MaundayThursday   time.Time `json:"MaundayThursday" yaml:"MaundayThursday" bson:"MaundayThursday"`
// 	GoodFriday        time.Time `json:"GoodFriday" yaml:"GoodFriday" bson:"GoodFriday"`
// 	EasterMonday      time.Time `json:"EasterMonday" yaml:"EasterMonday" bson:"EasterMonday"`
// 	LaborDay          time.Time `json:"LaborDay" yaml:"LaborDay" bson:"LaborDay"`
// 	EuropeDay         time.Time `json:"EuropeDay" yaml:"EuropeDay" bson:"EuropeDay"`
// 	AscensionThursday time.Time `json:"AscensionThursday" yaml:"AscensionThursday" bson:"AscensionThursday"`
//  WhitMonday        time.Time `json:"WhitMonday" yaml:"WhitMonday" bson:"WhitMonday"`
//  CorpusChristi     time.Time `json:"CorpusChristi" yaml:"CorpusChristi" bson:"CorpusChristi"`
//  GermanUnity       time.Time `json:"GermanUnity" yaml:"GermanUnity" bson:"GermanUnity"`
//  AlSaintsDay       time.Time `json:"AlSaintsDay" yaml:"AlSaintsDay" bson:"AlSaintsDay"`
//  ChristmasEve      time.Time `json:"ChristmasEve" yaml:"ChristmasEve" bson:"ChristmasEve"`
// 	ChristmasDay      time.Time `json:"ChristmasDay" yaml:"ChristmasDay" bson:"ChristmasDay"`
// 	ChristmasHoliday  time.Time `json:"ChristmasHoliday" yaml:"ChristmasHoliday" bson:"ChristmasHoliday"`
// }

// DEHolidays German holiday structure
type DEHolidays struct {
	Neujahrstag               time.Time `json:"Neujahrstag" yaml:"Neujahrstag" bson:"Neujahrstag"`                                           // New Years day
	Karfreitag                time.Time `json:"Karfreitag" yaml:"Karfreitag" bson:"Karfreitag"`                                              // Good Friday
	Ostermontag               time.Time `json:"Ostermontag" yaml:"Ostermontag" bson:"Ostermontag"`                                           // Easter Monday
	TagderArbeit              time.Time `json:"TagderArbeit" yaml:"TagderArbeit" bson:"TagderArbeit"`                                        // Labor day, May 1st
	ChristiHimmelfahrt        time.Time `json:"ChristiHimmelfahrt" yaml:"ChristiHimmelfahrt" bson:"ChristiHimmelfahrt"`                      // Ascension Day Easter Sunday + 39d
	Pfingstmontag             time.Time `json:"Pfingstmontag" yaml:"Pfingstmontag" bson:"Pfingstmontag"`                                     // Whit Monday Easter Sunday + 50d
	TagderDeutschenEinheit    time.Time `json:"TagderDeutschenEinheit" yaml:"TagderDeutschenEinheit" bson:"TagderDeutschenEinheit"`          // German Unity Day, October 3rd
	Weihnachtstag             time.Time `json:"Weihnachtstag" yaml:"Weihnachtstag" bson:"Weihnachtstag"`                                     // Christmas Day
	ZweiterWeihnachtsfeiertag time.Time `json:"ZweiterWeihnachtsfeiertag" yaml:"ZweiterWeihnachtsfeiertag" bson:"ZweiterWeihnachtsfeiertag"` // St Stephen's Day / Boxing Day December 26th
}

// NYSEHolidaysObserved US NYSE holiday structure
// Washington's Birthday: Though other institutions such as state and local
// governments and private businesses may use "president's day",
// it is Federal policy to always refer to holidays
// by the names designated in the law.
type NYSEHolidaysObserved struct {
	NewYearsDay         time.Time `json:"NewYearsDay" yaml:"NewYearsDay" bson:"NewYearsDay"`
	MartinLutherKing    time.Time `json:"MartinLutherKing" yaml:"MartinLutherKing" bson:"MartinLutherKing"`
	WashingtonsBirthday time.Time `json:"WashingtonsBirthday" yaml:"WashingtonsBirthday" bson:"WashingtonsBirthday"`
	GoodFriday          time.Time `json:"GoodFriday" yaml:"GoodFriday" bson:"GoodFriday"`
	MemorialDay         time.Time `json:"MemorialDay" yaml:"MemorialDay" bson:"MemorialDay"`
	IndependenceDay     time.Time `json:"IndependenceDay" yaml:"IndependenceDay" bson:"IndependenceDay"`
	LaborDay            time.Time `json:"LaborDay" yaml:"LaborDay" bson:"LaborDay"`
	ThanksgivingDay     time.Time `json:"ThanksgivingDay" yaml:"ThanksgivingDay" bson:"ThanksgivingDay"`
	ChristmasDay        time.Time `json:"ChristmasDay" yaml:"ChristmasDay" bson:"ChristmasDay"`
}

// USFederalHolidaysObserved US Federal holiday structure
//Washington's Birthday: Though other institutions such as state and local
//governments and private businesses may use 'presidents day',
//it is Federal policy to always refer to holidays
//by the names designated in the law.
type USFederalHolidaysObserved struct {
	NewYearsDay         time.Time `json:"NewYearsDay" yaml:"NewYearsDay" bson:"NewYearsDay"`
	MartinLutherKing    time.Time `json:"MartinLutherKing" yaml:"MartinLutherKing" bson:"MartinLutherKing"`
	WashingtonsBirthday time.Time `json:"WashingtonsBirthday" yaml:"WashingtonsBirthday" bson:"WashingtonsBirthday"`
	MemorialDay         time.Time `json:"MemorialDay" yaml:"MemorialDay" bson:"MemorialDay"`
	IndependenceDay     time.Time `json:"IndependenceDay" yaml:"IndependenceDay" bson:"IndependenceDay"`
	LaborDay            time.Time `json:"LaborDay" yaml:"LaborDay" bson:"LaborDay"`
	ColumbusDay         time.Time `json:"ColumbusDay" yaml:"ColumbusDay" bson:"ColumbusDay"`
	VeteransDay         time.Time `json:"VeteransDay" yaml:"VeteransDay" bson:"VeteransDay"`
	ThanksgivingDay     time.Time `json:"ThanksgivingDay" yaml:"ThanksgivingDay" bson:"ThanksgivingDay"`
	ChristmasDay        time.Time `json:"ChristmasDay" yaml:"ChristmasDay" bson:"ChristmasDay"`
}

// AustrailianHolidays Austrilian holiday structure
type AustrailianHolidays struct {
	NewYearsDay   time.Time `json:"NewYearsDay" yaml:"NewYearsDay" bson:"NewYearsDay"`
	Austrailiaday time.Time `json:"Austrailiaday" yaml:"Austrailiaday" bson:"Austrailiaday"`
	GoodFriday    time.Time `json:"GoodFriday" yaml:"GoodFriday" bson:"GoodFriday"`
	EasterMonday  time.Time `json:"EasterMonday" yaml:"EasterMonday" bson:"EasterMonday"`
	ANZACDday     time.Time `json:"ANZACDday" yaml:"ANZACDday" bson:"ANZACDday"`
	ChristmasDay  time.Time `json:"ChristmasDay" yaml:"ChristmasDay" bson:"ChristmasDay"`
}

// JapanBankHolidays Japanese Bank holiday structure
// note: Beginning in 2000, Japan implemented the Happy Monday System,
// which moved a number of national holidays to Monday in order to obtain a long weekend
type JapanBankHolidays struct {
	NewYearsDay           time.Time `json:"NewYearsDay" yaml:"NewYearsDay" bson:"NewYearsDay"`
	BankHoliday2          time.Time `json:"BankHoliday2" yaml:"BankHoliday2" bson:"BankHoliday2"`
	BankHoliday3          time.Time `json:"BankHoliday3" yaml:"BankHoliday3" bson:"BankHoliday3"`
	ComingOfAgeDay        time.Time `json:"ComingOfAgeDay" yaml:"ComingOfAgeDay" bson:"ComingOfAgeDay"`
	NationalFoundationDay time.Time `json:"NationalFoundationDay" yaml:"NationalFoundationDay" bson:"NationalFoundationDay"`
	VernalEquinoxDay      time.Time `json:"VernalEquinoxDay" yaml:"VernalEquinoxDay" bson:"VernalEquinoxDay"`
	ShowaDay              time.Time `json:"ShowaDay" yaml:"ShowaDay" bson:"ShowaDay"`
	ConstitutionDay       time.Time `json:"ConstitutionDay" yaml:"ConstitutionDay" bson:"ConstitutionDay"`
	GreeneryDay           time.Time `json:"GreeneryDay" yaml:"GreeneryDay" bson:"GreeneryDay"`
	ChildrensDay          time.Time `json:"ChildrensDay" yaml:"ChildrensDay" bson:"ChildrensDay"`
	MarineDay             time.Time `json:"MarineDay" yaml:"MarineDay" bson:"MarineDay"`
	MountainDay           time.Time `json:"MountainDay" yaml:"MountainDay" bson:"MountainDay"`
	RespectForTheAgedDay  time.Time `json:"RespectForTheAgedDay" yaml:"RespectForTheAgedDay" bson:"RespectForTheAgedDay"`
	AutumnalEquinoxDay    time.Time `json:"AutumnalEquinoxDay" yaml:"AutumnalEquinoxDay" bson:"AutumnalEquinoxDay"`
	HealthSportsDay       time.Time `json:"HealthSportsDay" yaml:"HealthSportsDay" bson:"HealthSportsDay"`
	CultureDay            time.Time `json:"CultureDay" yaml:"CultureDay" bson:"CultureDay"`
	LaborThanksgivingDay  time.Time `json:"LaborThanksgivingDay" yaml:"LaborThanksgivingDay" bson:"LaborThanksgivingDay"`
	EmperorsBirthday      time.Time `json:"EmperorsBirthday" yaml:"EmperorsBirthday" bson:"EmperorsBirthday"`
	NewYearsEve           time.Time `json:"NewYearsEve" yaml:"NewYearsEve" bson:"NewYearsEve"`
}

//Americas regional holiday structure
type Americas struct {
	USFederalHolidaysObserved USFederalHolidaysObserved `json:"USFederalHolidaysObserved" yaml:"USFederalHolidaysObserved" bson:"USFederalHolidaysObserved"`
	NYSEHolidaysObserved      NYSEHolidaysObserved      `json:"NYSEHolidaysObserved" yaml:"NYSEHolidaysObserved" bson:"NYSEHolidaysObserved"`
}

// Europe regional holiday structure
type Europe struct {
	DEHolidays         DEHolidays         `json:"DEHolidays" yaml:"DEHolidays" bson:"DEHolidays"`
	ECBTarget2Holidays ECBTarget2Holidays `json:"ECBTarget2Holidays" yaml:"ECBTarget2Holidays" bson:"ECBTarget2Holidays"`
	NLHolidays         NLHolidays         `json:"NLHolidays" yaml:"NLHolidays" bson:"NLHolidays"`
	UKHolidays         UKHolidays         `json:"UKHolidays" yaml:"UKHolidays" bson:"UKHolidays"`
}

// AsiaPacific regional holidays
type AsiaPacific struct {
	AustrailianHolidays AustrailianHolidays `json:"AustrailianHolidays" yaml:"AustrailianHolidays" bson:"AustrailianHolidays"`
	JapanBankHolidays   JapanBankHolidays   `json:"JapanBankHolidays" yaml:"JapanBankHolidays" bson:"JapanBankHolidays"`
}

// Holidays Master holidays structure
type Holidays struct {
	Year        int         `json:"Year" yaml:"Year" bson:"Year"`
	Americas    Americas    `json:"Americas" yaml:"Americas" bson:"Americas"`
	Europe      Europe      `json:"Europe" yaml:"Europe" bson:"Europe"`
	AsiaPacific AsiaPacific `json:"AsiaPacific" yaml:"AsiaPacific" bson:"AsiaPacific"`
}

// IsWeekend Function to determine if the date falls on a weekend (SAT or SUN).
func IsWeekend(date time.Time) bool {
	day := date.Weekday()
	return day == time.Saturday || day == time.Sunday
}

// return3rdMonday function return the 3rd Monday of the Month
func return3rdMonday(yyyy int, mm time.Month) time.Time {
	date := time.Date(yyyy, mm, 1, 0, 0, 0, 0, time.UTC)
	switch date.Weekday() {
	case time.Tuesday:
		return date.AddDate(0, 0, 20)
	case time.Wednesday:
		return date.AddDate(0, 0, 19)
	case time.Thursday:
		return date.AddDate(0, 0, 18)
	case time.Friday:
		return date.AddDate(0, 0, 17)
	case time.Saturday:
		return date.AddDate(0, 0, 16)
	case time.Sunday:
		return date.AddDate(0, 0, 16)
	case time.Monday:
		return date.AddDate(0, 0, 14)
	}
	//should never get here but go compiler wants a return after the case
	return date
}

// returnMonthEnd reports the ending day of the month in t
func returnMonthEnd(date time.Time) time.Time {
	return time.Date(date.Year(), date.Month()+1, 0, date.Hour(), date.Minute(), date.Second(), date.Nanosecond(), date.Location())
}

func returnFirstMonday(yyyy int, mm time.Month) time.Time {
	date := time.Date(yyyy, mm, 1, 0, 0, 0, 0, time.UTC)
	switch date.Weekday() {
	case time.Monday:
		return date
	case time.Tuesday:
		return date.AddDate(0, 0, 6)
	case time.Wednesday:
		return date.AddDate(0, 0, 5)
	case time.Thursday:
		return date.AddDate(0, 0, 4)
	case time.Friday:
		return date.AddDate(0, 0, 3)
	case time.Saturday:
		return date.AddDate(0, 0, 2)
	case time.Sunday:
		return date.AddDate(0, 0, 1)
	}
	//should never get here but go compiler wants a return after the case
	return date
}

// returnSecondMonday used to calculate Canadian Thanksgiving Day
// amd US Columbus Day
func returnSecondMonday(yyyy int, mm time.Month) time.Time {
	return (returnFirstMonday(yyyy, mm)).AddDate(0, 0, 7)
}

// returnLastMonday Return the last Monday in the month
func returnLastMonday(yyyy int, mm time.Month) time.Time {
	date := returnMonthEnd(time.Date(yyyy, mm, 1, 0, 0, 0, 0, time.UTC))
	switch date.Weekday() {
	case time.Monday:
		return date
	case time.Tuesday:
		return date.AddDate(0, 0, -1)
	case time.Wednesday:
		return date.AddDate(0, 0, -2)
	case time.Thursday:
		return date.AddDate(0, 0, -3)
	case time.Friday:
		return date.AddDate(0, 0, -4)
	case time.Saturday:
		return date.AddDate(0, 0, -5)
	case time.Sunday:
		return date.AddDate(0, 0, -6)
	}
	//should never get here but go compiler wants a return after the case
	return date
}

func returnFourthThursday(yyyy int, mm time.Month) time.Time {
	date := time.Date(yyyy, mm, 1, 0, 0, 0, 0, time.UTC)
	switch date.Weekday() {
	case time.Monday:
		return date.AddDate(0, 0, 24)
	case time.Tuesday:
		return date.AddDate(0, 0, 23)
	case time.Wednesday:
		return date.AddDate(0, 0, 22)
	case time.Thursday:
		return date.AddDate(0, 0, 21)
	case time.Friday:
		return date.AddDate(0, 0, 20)
	case time.Saturday:
		return date.AddDate(0, 0, 19)
	case time.Sunday:
		return date.AddDate(0, 0, 18)
	}
	//should never get here but go compiler wants a return after the case
	return date
}

func returnObservableChristmas(yyyy int) time.Time {
	date := time.Date(yyyy, time.December, 25, 0, 0, 0, 0, time.UTC)
	switch date.Weekday() {
	case time.Saturday:
		return date.AddDate(0, 0, -1)
	case time.Sunday:
		return date.AddDate(0, 0, 1)
	}
	return date
}

func returnObservableUSVeterandDay(yyyy int) time.Time {
	date := time.Date(yyyy, time.November, 11, 0, 0, 0, 0, time.UTC)
	switch date.Weekday() {
	case time.Saturday:
		return date.AddDate(0, 0, -1)
	case time.Sunday:
		return date.AddDate(0, 0, 1)
	}
	return date
}

// calculateGregorianEaster Calculate Gregorian Calendar Easter date
func calculateGregorianEaster(year int) time.Time {
	// This function uses the algorithm invented by the mathematician
	// Carl Friedrich Gauss in 1800 to calculate the date of Easter in a given year
	// returns day, month, year as integers
	//
	// don't go below start of Gregorian calendar
	const firstgregorianyear int = 1583
	// don't go above the year where integer calculations will start to fail
	const lastgregorianyear int = 4099
	yyyy := year
	if yyyy < firstgregorianyear {
		yyyy = firstgregorianyear
	} else if year > lastgregorianyear {
		yyyy = lastgregorianyear
	}
	// start month off in March since this is the earliest it can be.
	mm := 3
	// determine the Golden number
	goldennumber := ((yyyy % 19) + 1)
	// determine the century number
	centurynumber := yyyy/100 + 1
	// correct for the years that are not leap years
	xx := (3*centurynumber)/4 - 12
	// moon correction
	yy := (8*centurynumber+5)/25 - 5
	// find Sunday
	zz := (5*yyyy)/4 - xx - 10
	// determine epoch
	// age of moon on January 1st of the year
	// that follows a cycle of every 19 years
	ee := (11*goldennumber + 20 + yy - xx) % 30
	if ee == 24 {
		ee++
	}
	if (ee == 25) && (goldennumber > 11) {
		ee++
	}
	// get the full moon
	moon := (44 - ee)
	if moon < 21 {
		moon += 30
	}
	// up to Sunday
	dd := (moon + 7) - ((zz + moon) % 7)
	// possibly up a month in easter_date
	if dd > 31 {
		dd -= 31
		mm = 4
	}
	if mm == 3 {
		return time.Date(yyyy, time.March, dd, 0, 0, 0, 0, time.UTC)
	}
	// else return April since it is not March
	return time.Date(yyyy, time.April, dd, 0, 0, 0, 0, time.UTC)
}

//calculateGregorianGoodFriday two days before easter.
func calculateGregorianGoodFriday(year int) time.Time {
	//calculate easter and subtract two days
	return (calculateGregorianEaster(year)).AddDate(0, 0, -2)
}

//calculateGregorianEasterMonday 1 days after easter.
func calculateGregorianEasterMonday(year int) time.Time {
	//calculate easter and add 1 day
	return (calculateGregorianEaster(year)).AddDate(0, 0, 1)
}

//calculateGregorianAscension 40 days after easter.
func calculateGregorianAscension(year int) time.Time {
	//calculate easter and add 40 days
	return (calculateGregorianEaster(year)).AddDate(0, 0, 40)
}

//calculateGregorianPentecost 50 days after easter.
func calculateGregorianPentecost(year int) time.Time {
	//calculate easter and add 40 days
	return (calculateGregorianEaster(year)).AddDate(0, 0, 50)
}

// inBetween : checks if i is between the min and the max returns boolean
func inBetween(i, min, max int) bool {
	if (i >= min) && (i <= max) {
		return true
	}

	return false

}

func createUniqueFileString(n int) string {
	b := make([]byte, 5)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return fmt.Sprintf("%X", b)
}

func setObservableHoliday(inDate time.Time, moveToMonday bool) time.Time {
	if inDate.Weekday() == time.Saturday || inDate.Weekday() == time.Sunday {
		if inDate.Weekday() == time.Saturday {
			return inDate.AddDate(0, 0, 2)
		}
		return inDate.AddDate(0, 0, 1)
	}
	return inDate
}

func setHolidayDE(yyyy int, h *Holidays) {
	//Neujahrstag = New Years day
	h.Europe.DEHolidays.Neujahrstag = setObservableHoliday(time.Date(yyyy, time.January, 1, 0, 0, 0, 0, time.UTC), moveToMonday)
	//Karfreitag  = Good Friday
	h.Europe.DEHolidays.Karfreitag = calculateGregorianGoodFriday(yyyy)
	//Ostermontag = Easter Monday
	h.Europe.DEHolidays.Ostermontag = calculateGregorianEasterMonday(yyyy)
	//TagderArbeit = Labor day, May 1st
	h.Europe.DEHolidays.TagderArbeit = time.Date(yyyy, time.May, 1, 0, 0, 0, 0, time.UTC)
	//ChristiHimmelfahrt  = Ascension Day Easter Sunday + 39 days
	h.Europe.DEHolidays.ChristiHimmelfahrt = calculateGregorianAscension(yyyy)
	//Pfingstmontag  = Whit Monday Easter Sunday + 50d
	h.Europe.DEHolidays.Pfingstmontag = calculateGregorianPentecost(yyyy)
	//TagderDeutschenEinheit = German Unity Day, October 3rd
	h.Europe.DEHolidays.TagderDeutschenEinheit = time.Date(yyyy, time.October, 3, 0, 0, 0, 0, time.UTC)
	//Weihnachtstag = Christmas Day
	h.Europe.DEHolidays.Weihnachtstag = time.Date(yyyy, time.December, 25, 0, 0, 0, 0, time.UTC)
	//ZweiterWeihnachtsfeiertag = St Stephen's Day / Boxing Day December 26th
	h.Europe.DEHolidays.ZweiterWeihnachtsfeiertag = time.Date(yyyy, time.December, 26, 0, 0, 0, 0, time.UTC)
	fmt.Println("setHolidayDE completed...")
	waitGroup.Done()
}

// setHolidaysNL : Netherland holidays
// Bevrijdingsdag is a holiday every 5 years
func setHolidaysNL(yyyy int, h *Holidays) {
	// Nieuwjaar       = NewYear
	h.Europe.NLHolidays.Nieuwjaardag = setObservableHoliday(time.Date(yyyy, time.January, 1, 0, 0, 0, 0, time.UTC), moveToMonday)
	// GoedeVrijdag    = GoodFriday - not an official public holiday,
	// though most schools are closed. Some public offices are closed
	//  but most commercial concerns like banks and stores are open.
	h.Europe.NLHolidays.Goedevrijdag = calculateGregorianGoodFriday(yyyy)
	// PaasMaandag     = EasterMonday
	h.Europe.NLHolidays.Paasmaandag = calculateGregorianEasterMonday(yyyy)
	// KoningsDag      = KoningsDag April 27th. If Sunday then observed Saturday
	if time.Date(yyyy, time.April, 27, 0, 0, 0, 0, time.UTC).Weekday() == time.Sunday {
		h.Europe.NLHolidays.Koningsdag = time.Date(yyyy, time.April, 26, 0, 0, 0, 0, time.UTC)
	} else {
		h.Europe.NLHolidays.Koningsdag = time.Date(yyyy, time.April, 27, 0, 0, 0, 0, time.UTC)
	}
	// BevrijdingsDag  = May, 5
	h.Europe.NLHolidays.Bevrijdingsdag = time.Date(yyyy, time.May, 5, 0, 0, 0, 0, time.UTC)
	// Hemelvaart      = Ascension, 40 days after Easter
	h.Europe.NLHolidays.Hemelvaart = calculateGregorianAscension(yyyy)
	// PinksterMaandag = Whit Sunday/Pentecost  50 days after Easter
	h.Europe.NLHolidays.Pinkstermaandag = calculateGregorianPentecost(yyyy)
	// EersteKerstdag  = Christmas
	h.Europe.NLHolidays.Eerstekerstdag = time.Date(yyyy, time.December, 25, 0, 0, 0, 0, time.UTC)
	// TweedeKerstdag  = Christmas2
	h.Europe.NLHolidays.Tweedekerstdag = time.Date(yyyy, time.December, 26, 0, 0, 0, 0, time.UTC)
	fmt.Println("setHolidaysNL completed...")
	waitGroup.Done()
}

// setHolidaysUK : set the UK Bank holidays
func setHolidaysUK(yyyy int, h *Holidays) {
	// New Years day Observed
	h.Europe.UKHolidays.NewYearsDay = setObservableHoliday(time.Date(yyyy, time.January, 1, 0, 0, 0, 0, time.UTC), moveToMonday)
	h.Europe.UKHolidays.GoodFriday = calculateGregorianGoodFriday(yyyy)
	h.Europe.UKHolidays.EarlyMay = returnFirstMonday(yyyy, time.May)
	h.Europe.UKHolidays.SpringHoliday = returnLastMonday(yyyy, time.May)
	switch true {
	case time.Date(yyyy, time.December, 26, 0, 0, 0, 0, time.UTC).Weekday() == time.Saturday:
		//if Christmas is Friday then Boxing day is Observed Monday
		h.Europe.UKHolidays.ChristmasDay = time.Date(yyyy, time.December, 25, 0, 0, 0, 0, time.UTC)
		h.Europe.UKHolidays.BoxingDay = time.Date(yyyy, time.December, 28, 0, 0, 0, 0, time.UTC)
	case time.Date(yyyy, time.December, 25, 0, 0, 0, 0, time.UTC).Weekday() == time.Saturday:
		// if Christmas and Boxing day fall on the weekend then Monday and Tuesday are observed
		h.Europe.UKHolidays.ChristmasDay = time.Date(yyyy, time.December, 27, 0, 0, 0, 0, time.UTC)
		h.Europe.UKHolidays.BoxingDay = time.Date(yyyy, time.December, 28, 0, 0, 0, 0, time.UTC)
	case time.Date(yyyy, time.December, 25, 0, 0, 0, 0, time.UTC).Weekday() == time.Sunday:
		// if Christmas is on Sunday then it is observed Monday and Boxing Day on Tuesday
		h.Europe.UKHolidays.ChristmasDay = time.Date(yyyy, time.December, 26, 0, 0, 0, 0, time.UTC)
		h.Europe.UKHolidays.BoxingDay = time.Date(yyyy, time.December, 27, 0, 0, 0, 0, time.UTC)
	default:
		h.Europe.UKHolidays.ChristmasDay = time.Date(yyyy, time.December, 25, 0, 0, 0, 0, time.UTC)
		h.Europe.UKHolidays.BoxingDay = time.Date(yyyy, time.December, 26, 0, 0, 0, 0, time.UTC)
	}
	fmt.Println("setHolidaysUK completed...")
	waitGroup.Done()
}

// setHolidaysECB : set the ECB holidays
func setHolidaysECB(yyyy int, h *Holidays) {
	// New Years day
	h.Europe.ECBTarget2Holidays.NewYearsDay = setObservableHoliday(time.Date(yyyy, time.January, 1, 0, 0, 0, 0, time.UTC), moveToMonday)
	h.Europe.ECBTarget2Holidays.GoodFriday = calculateGregorianGoodFriday(yyyy)
	h.Europe.ECBTarget2Holidays.EasterMonday = calculateGregorianEasterMonday(yyyy)
	h.Europe.ECBTarget2Holidays.LaborDay = time.Date(yyyy, time.May, 1, 0, 0, 0, 0, time.UTC)
	switch true {
	case time.Date(yyyy, time.December, 26, 0, 0, 0, 0, time.UTC).Weekday() == time.Saturday:
		//if Christmas is Friday then Christmas Holiday is Observed Monday
		h.Europe.ECBTarget2Holidays.ChristmasDay = time.Date(yyyy, time.December, 25, 0, 0, 0, 0, time.UTC)
		h.Europe.ECBTarget2Holidays.ChristmasHoliday = time.Date(yyyy, time.December, 28, 0, 0, 0, 0, time.UTC)
	case time.Date(yyyy, time.December, 25, 0, 0, 0, 0, time.UTC).Weekday() == time.Saturday:
		// if Christmas and Christmas Holiday fall on the weekend then Monday and Tuesday are observed
		h.Europe.ECBTarget2Holidays.ChristmasDay = time.Date(yyyy, time.December, 27, 0, 0, 0, 0, time.UTC)
		h.Europe.ECBTarget2Holidays.ChristmasHoliday = time.Date(yyyy, time.December, 28, 0, 0, 0, 0, time.UTC)
	case time.Date(yyyy, time.December, 25, 0, 0, 0, 0, time.UTC).Weekday() == time.Sunday:
		// if Christmas is on Sunday then it is observed Monday and Christmas Holiday on Tuesday
		h.Europe.ECBTarget2Holidays.ChristmasDay = time.Date(yyyy, time.December, 26, 0, 0, 0, 0, time.UTC)
		h.Europe.ECBTarget2Holidays.ChristmasHoliday = time.Date(yyyy, time.December, 27, 0, 0, 0, 0, time.UTC)
	default:
		h.Europe.ECBTarget2Holidays.ChristmasDay = time.Date(yyyy, time.December, 25, 0, 0, 0, 0, time.UTC)
		h.Europe.ECBTarget2Holidays.ChristmasHoliday = time.Date(yyyy, time.December, 26, 0, 0, 0, 0, time.UTC)
	}
	fmt.Println("setHolidaysECB completed...")
	waitGroup.Done()
}

// setHolidaysAustrailia : set the Austrailian Holidays -- assuming Australian Capital Territory rules
func setHolidaysAustrailia(yyyy int, h *Holidays) {
	// New Years day, if
	h.AsiaPacific.AustrailianHolidays.NewYearsDay = setObservableHoliday(time.Date(yyyy, time.January, 1, 0, 0, 0, 0, time.UTC), moveToMonday)
	h.AsiaPacific.AustrailianHolidays.Austrailiaday = setObservableHoliday(time.Date(yyyy, time.January, 26, 0, 0, 0, 0, time.UTC), moveToMonday)
	h.AsiaPacific.AustrailianHolidays.GoodFriday = calculateGregorianGoodFriday(yyyy)
	h.AsiaPacific.AustrailianHolidays.EasterMonday = calculateGregorianEasterMonday(yyyy)
	h.AsiaPacific.AustrailianHolidays.ANZACDday = setObservableHoliday(time.Date(yyyy, time.April, 25, 0, 0, 0, 0, time.UTC), moveToMonday)
	h.AsiaPacific.AustrailianHolidays.ChristmasDay = setObservableHoliday(time.Date(yyyy, time.December, 25, 0, 0, 0, 0, time.UTC), moveToMonday)
	fmt.Println("setHolidaysAustrailia completed...")
	waitGroup.Done()
}

// setHolidaysJapan set the japanese holidays...
func setHolidaysJapan(yyyy int, h *Holidays) {
	// don't go below start of Gregorian calendar
	if time.Date(yyyy, time.January, 1, 0, 0, 0, 0, time.UTC).Weekday() == time.Sunday {
		//Move new years to Monday
		h.AsiaPacific.JapanBankHolidays.NewYearsDay = time.Date(yyyy, time.January, 2, 0, 0, 0, 0, time.UTC)
		// bank holiday 2 and new years day share the same day
		h.AsiaPacific.JapanBankHolidays.BankHoliday2 = time.Date(yyyy, time.January, 2, 0, 0, 0, 0, time.UTC)
		h.AsiaPacific.JapanBankHolidays.BankHoliday3 = time.Date(yyyy, time.January, 3, 0, 0, 0, 0, time.UTC)
	} else {
		// if New Years Day falls on a Saturday then New Years Eve is a holiday otherwise it is nil
		if time.Date(yyyy, time.January, 1, 0, 0, 0, 0, time.UTC).Weekday() == time.Saturday {
			h.AsiaPacific.JapanBankHolidays.NewYearsEve = time.Date(yyyy-1, time.December, 31, 0, 0, 0, 0, time.UTC)
		} else {
			h.AsiaPacific.JapanBankHolidays.NewYearsEve = time.Time{}
		}
		h.AsiaPacific.JapanBankHolidays.NewYearsDay = time.Date(yyyy, time.January, 1, 0, 0, 0, 0, time.UTC)
		h.AsiaPacific.JapanBankHolidays.BankHoliday2 = time.Date(yyyy, time.January, 2, 0, 0, 0, 0, time.UTC)
		h.AsiaPacific.JapanBankHolidays.BankHoliday3 = time.Date(yyyy, time.January, 3, 0, 0, 0, 0, time.UTC)
	}
	// set coming of the aged date, 2nd Monday in January
	h.AsiaPacific.JapanBankHolidays.ComingOfAgeDay = returnSecondMonday(yyyy, time.January)
	// National Foundation Day is February 11th
	h.AsiaPacific.JapanBankHolidays.NationalFoundationDay = time.Date(yyyy, time.February, 3, 0, 0, 0, 0, time.UTC)
	// Vernal Equinox Day is calculated to fall on the Spring Equinox
	h.AsiaPacific.JapanBankHolidays.VernalEquinoxDay = time.Date(yyyy, time.March, julian.JDToTime(solstice.March(yyyy)).Day(), 0, 0, 0, 0, time.UTC)
	// Showa Day is April 29th, Part of Golden Week
	h.AsiaPacific.JapanBankHolidays.ShowaDay = time.Date(yyyy, time.April, 29, 0, 0, 0, 0, time.UTC)
	// Constitution Memorial Day is May 3rd, Part of Golden week
	h.AsiaPacific.JapanBankHolidays.ConstitutionDay = time.Date(yyyy, time.May, 3, 0, 0, 0, 0, time.UTC)
	// Greenery (Arbor) Day is May 4th, Part of Golden week
	h.AsiaPacific.JapanBankHolidays.GreeneryDay = time.Date(yyyy, time.May, 4, 0, 0, 0, 0, time.UTC)
	// Childrens Day is May 5th, Part of Golden week
	h.AsiaPacific.JapanBankHolidays.ChildrensDay = time.Date(yyyy, time.May, 4, 0, 0, 0, 0, time.UTC)
	// Marine day. First offical in 1996. Since 2003 this holiday is now the third Monday in July
	if yyyy >= 1996 && yyyy < 2003 {
		h.AsiaPacific.JapanBankHolidays.MarineDay = time.Date(yyyy, time.July, 20, 0, 0, 0, 0, time.UTC)
	} else {
		h.AsiaPacific.JapanBankHolidays.MarineDay = return3rdMonday(yyyy, time.July)
	}
	// Mountain Day is August 11th (or 12th if Sunday)
	if time.Date(yyyy, time.August, 11, 0, 0, 0, 0, time.UTC).Weekday() == time.Sunday {
		h.AsiaPacific.JapanBankHolidays.MountainDay = time.Date(yyyy, time.August, 12, 0, 0, 0, 0, time.UTC)
	} else {
		h.AsiaPacific.JapanBankHolidays.MountainDay = time.Date(yyyy, time.August, 11, 0, 0, 0, 0, time.UTC)
	}
	// Respect For The Aged Day is the 3rd Monday in September
	h.AsiaPacific.JapanBankHolidays.RespectForTheAgedDay = return3rdMonday(yyyy, time.September)
	// Autumnal Equinox Day  is calculated to fall on the Fall Equinox
	h.AsiaPacific.JapanBankHolidays.AutumnalEquinoxDay = time.Date(yyyy, time.March, julian.JDToTime(solstice.September(yyyy)).Day(), 0, 0, 0, 0, time.UTC)
	// Health Sports Day is the 2nd Monday in October
	h.AsiaPacific.JapanBankHolidays.HealthSportsDay = returnSecondMonday(yyyy, time.October)
	// Culture Day is November 3rd or the fourth if the 3rd is a Sunday
	if time.Date(yyyy, time.November, 3, 0, 0, 0, 0, time.UTC).Weekday() == time.Sunday {
		h.AsiaPacific.JapanBankHolidays.CultureDay = time.Date(yyyy, time.November, 4, 0, 0, 0, 0, time.UTC)
	} else {
		h.AsiaPacific.JapanBankHolidays.CultureDay = time.Date(yyyy, time.November, 3, 0, 0, 0, 0, time.UTC)
	}
	// Labor Thanksgiving Day is November 23rd or 24th if the 23rd is a Sunday
	if time.Date(yyyy, time.November, 23, 0, 0, 0, 0, time.UTC).Weekday() == time.Sunday {
		h.AsiaPacific.JapanBankHolidays.LaborThanksgivingDay = time.Date(yyyy, time.November, 24, 0, 0, 0, 0, time.UTC)
	} else {
		h.AsiaPacific.JapanBankHolidays.LaborThanksgivingDay = time.Date(yyyy, time.November, 23, 0, 0, 0, 0, time.UTC)
	}
	// Emperors Birthday : effective 1989 is now December 23rd.
	if yyyy >= 1989 {
		h.AsiaPacific.JapanBankHolidays.EmperorsBirthday = time.Date(yyyy, time.December, 23, 0, 0, 0, 0, time.UTC)
	}
	fmt.Println("setHolidaysJapan completed...")
	waitGroup.Done()
}

// setHolidaysUS : set US holidays...
func setHolidaysUS(yyyy int, h *Holidays) {
	// Calculate the US Federal Holidays
	h.Americas.USFederalHolidaysObserved.NewYearsDay = setObservableHoliday(time.Date(yyyy, time.January, 1, 0, 0, 0, 0, time.UTC), moveToMonday)
	// Martin Luther King is the third Monday in January
	h.Americas.USFederalHolidaysObserved.MartinLutherKing = return3rdMonday(yyyy, time.January)
	// Washington's Birthday is the third Monday in February
	h.Americas.USFederalHolidaysObserved.WashingtonsBirthday = return3rdMonday(yyyy, time.February)
	// Memorial Day is the last Monday in May
	h.Americas.USFederalHolidaysObserved.MemorialDay = returnLastMonday(yyyy, time.May)
	// Independence Day is July 4th or nearest Monday
	h.Americas.USFederalHolidaysObserved.IndependenceDay = setObservableHoliday(time.Date(yyyy, time.July, 4, 0, 0, 0, 0, time.UTC), moveToMonday)
	// Labor Day is the first Monday in September
	h.Americas.USFederalHolidaysObserved.LaborDay = returnFirstMonday(yyyy, time.September)
	// Columbus Day is the Second Monday in October
	h.Americas.USFederalHolidaysObserved.ColumbusDay = returnSecondMonday(yyyy, time.October)
	// Veterans Day is Novemeber 11th
	h.Americas.USFederalHolidaysObserved.VeteransDay = returnObservableUSVeterandDay(yyyy)
	//Thanksgiving Day is the Fourth Thursday of the month
	h.Americas.USFederalHolidaysObserved.ThanksgivingDay = returnFourthThursday(yyyy, time.November)
	//Christmas Day if falls on a Saturday is observed Friday or if it falls on a Sunday is observed Monday
	h.Americas.USFederalHolidaysObserved.ChristmasDay = returnObservableChristmas(yyyy)
	// Copy US Federal Holdiays and then calculate Good Friday
	h.Americas.NYSEHolidaysObserved.NewYearsDay = h.Americas.USFederalHolidaysObserved.NewYearsDay
	h.Americas.NYSEHolidaysObserved.MartinLutherKing = h.Americas.USFederalHolidaysObserved.MartinLutherKing
	h.Americas.NYSEHolidaysObserved.WashingtonsBirthday = h.Americas.USFederalHolidaysObserved.WashingtonsBirthday
	// Good Friday is two days before Gregorian/Western Easter
	h.Americas.NYSEHolidaysObserved.GoodFriday = calculateGregorianGoodFriday(yyyy)
	h.Americas.NYSEHolidaysObserved.MemorialDay = h.Americas.USFederalHolidaysObserved.MemorialDay
	h.Americas.NYSEHolidaysObserved.IndependenceDay = h.Americas.USFederalHolidaysObserved.IndependenceDay
	h.Americas.NYSEHolidaysObserved.LaborDay = h.Americas.USFederalHolidaysObserved.LaborDay
	h.Americas.NYSEHolidaysObserved.ThanksgivingDay = h.Americas.USFederalHolidaysObserved.ThanksgivingDay
	h.Americas.NYSEHolidaysObserved.ChristmasDay = h.Americas.USFederalHolidaysObserved.ChristmasDay
	fmt.Println("setHolidaysUS completed...")
	waitGroup.Done()
}

//Set the holiday year attribute in the structure
func setHolidayYear(yyyy int, h *Holidays) {
	// don't go below start of Gregorian calendar
	const firstgregorianyear int = 1583
	// don't go above the year where integer calculations will start to fail
	const lastgregorianyear int = 4099
	if inBetween(yyyy, firstgregorianyear, lastgregorianyear) {
		h.Year = yyyy
	}
	fmt.Println("setHolidayYear completed...")
	waitGroup.Done()
}

// checkIfFileOrDirectoryExists Check to see if the provided file, path or
// file and path combination exists.
func checkIfFileOrDirectoryExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

// Main function...
func main() {
	outputFilePath := "X:\\go\\output"
	var currentDirectory string
	var err error
	var pathSeparator string
	// set the proper path separator for the compiled file system
	if runtime.GOOS == "windows" {
		pathSeparator = "\\"
	} else {
		// then use the linux path separator
		pathSeparator = "/"
	}
	// check to see if the supplied directory exists

	fileOrDirectoryExists, _ := checkIfFileOrDirectoryExists(outputFilePath)
	if fileOrDirectoryExists {
		currentDirectory = outputFilePath
	} else {
		// get the current working directory where the executable is running
		currentDirectory, err = os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
	}

	// create unique file name suffix for the output files and the names
	// use the date, time and a random number to insure unique file names
	// for the execution of the program
	s := fmt.Sprintf((time.Now()).Format("20060102150405")) + "-" + strings.ToLower(createUniqueFileString(5))
	jsonfile := currentDirectory + pathSeparator + "holidays-" + s + ".json"
	fmt.Println("Creating JSON File: ", jsonfile)
	xmlfile := currentDirectory + pathSeparator + "holidays-" + s + ".xml"
	fmt.Println("Creating XML File: ", xmlfile)
	yamlfile := currentDirectory + pathSeparator + "holidays-" + s + ".yaml"
	fmt.Println("Creating YAML File: ", yamlfile)
	bsonfile := currentDirectory + pathSeparator + "holidays-" + s + ".bson"
	fmt.Println("Creating BSON File: ", bsonfile)
	// get the current UTC time as of now... a.k.a. current date
	now := (time.Now()).UTC()
	processYear := now.Year()
	// create holidays struct and assign to h
	h := Holidays{}
	// create a pointer to the structure holidays
	pointer2h := &h
	/////////////////////////////////////////////////////////////////////////////////////
	// Processs each country as a separate thread since they are not dependent
	// upon each other nor do they share any variables other than being part of the
	// overal structure.

	// use the waitGroup to track the code finishing.
	{
		waitGroup.Add(8)
		go setHolidayYear(processYear, pointer2h)
		go setHolidaysNL(processYear, pointer2h)
		go setHolidayDE(processYear, pointer2h)
		go setHolidaysUS(processYear, pointer2h)
		go setHolidaysJapan(processYear, pointer2h)
		go setHolidaysAustrailia(processYear, pointer2h)
		go setHolidaysUK(processYear, pointer2h)
		go setHolidaysECB(processYear, pointer2h)
		waitGroup.Wait()
	}
	//
	//marshal out the structure as json...
	//fmt.Printf("\nmarshal out the structure as json...\n")
	jsonout, jsonMarshalError := json.Marshal(h)
	if jsonMarshalError != nil {
		log.Println(jsonMarshalError)
	}
	//fmt.Println(string(jsonout))
	jsonioerr := ioutil.WriteFile(jsonfile, jsonout, outputFilePermission)
	if jsonioerr != nil {
		log.Println(jsonioerr)
	}
	//
	//marshal out the structure as yaml..
	//fmt.Printf("\nmarshal out the structure as yaml...\n")
	yamlout, yamlMarshalError := yaml.Marshal(h)
	if yamlMarshalError != nil {
		log.Println(yamlMarshalError)
	}
	//fmt.Println(string(yamlout))
	yamlioerr := ioutil.WriteFile(yamlfile, yamlout, outputFilePermission)
	if yamlioerr != nil {
		log.Println(yamlioerr)
	}
	//
	//marshal out the structure as bson..
	//fmt.Printf("\nmarshal out the structure as bson...\n")
	bsonout, bsonMarshalError := bson.Marshal(h)
	if bsonMarshalError != nil {
		log.Println(bsonMarshalError)
	}
	//fmt.Printf("%q", bsonout)
	bsonioerr := ioutil.WriteFile(bsonfile, bsonout, outputFilePermission)
	if bsonioerr != nil {
		log.Println(bsonioerr)
	}

	//fmt.Printf("\n\nmarshal out the structure as XML...\n")
	xmlout, xmlMarshalError := xml.Marshal(h)
	if xmlMarshalError != nil {
		log.Println(xmlMarshalError)
	}
	//fmt.Printf("%q", xmlout)
	xmlioerr := ioutil.WriteFile(xmlfile, xmlout, outputFilePermission)
	if xmlioerr != nil {
		log.Println(xmlioerr)
	}
	//
	fmt.Println("\n\n\nEnd of func (main) ...")
}
