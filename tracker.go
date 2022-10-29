package santatracker

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"time"
)

const PresentsDeliveredOverWater float64 = 0.3
const PresentsDeliveredOverLand float64 = 1 - PresentsDeliveredOverWater

func init() {
	log.SetFlags(0)
	log.SetPrefix("[Tracker] ")
}

type Tracker struct {
	currentIndex int
	locationData *LocationData
	timeOffset   int64
}

func (t *Tracker) updateIndex() {
	now := t.getAdjustedTime()

	// Hasn't set off yet
	if now < t.locationData.Destinations[0].Departure {
		t.currentIndex = 0
		return
	}

	// At the last location
	lastDestinationIndex := len(t.locationData.Destinations) - 1
	if now > t.locationData.Destinations[lastDestinationIndex].Arrival {
		t.currentIndex = lastDestinationIndex
		return
	}

	/*
		If we're not quite at the next destination, we're
		still at the previous one (hence, idx -1)
	*/
	for idx, dest := range t.locationData.Destinations {
		if now < dest.Arrival {
			t.currentIndex = idx - 1
			break
		}
	}
}

func (t *Tracker) getAdjustedTime() int64 {
	return time.Now().UnixMilli() + t.timeOffset
}

func (t *Tracker) next() Destination {
	lastDestinationIndex := len(t.locationData.Destinations) - 1
	if t.currentIndex == lastDestinationIndex {
		return t.locationData.Destinations[lastDestinationIndex]
	}
	return t.locationData.Destinations[t.currentIndex+1]
}

func (t *Tracker) prev() Destination {
	if t.currentIndex == 0 {
		return t.locationData.Destinations[0]
	}
	return t.locationData.Destinations[t.currentIndex-1]
}

// TODO: Update this to take into account presents delivered over water and land
func (t *Tracker) calculatePresentsDelivered(now int64, current Destination, next Destination) int64 {

	// How long we've been at this destination
	fElapsedMs := float64(now - current.Arrival)
	fTotalToDeliver := float64(next.PresentsDelivered - current.PresentsDelivered)
	if fTotalToDeliver == 0 || fElapsedMs <= 0 {
		return current.PresentsDelivered
	}

	fDurationMs := float64(next.Arrival - current.Arrival)
	fTotalDeliveredSoFar := math.Floor((fElapsedMs / fDurationMs) * fTotalToDeliver)

	log.Printf("[Cur. Period] Elap.: %.f, Dura.: %.f, ToDel.: %.f, DelSoFa.: %.f%%", fElapsedMs, fDurationMs, fTotalToDeliver, (fTotalDeliveredSoFar/fTotalToDeliver)*100)

	return int64(fTotalDeliveredSoFar + float64(current.PresentsDelivered))
}

func (t *Tracker) GetCurrentLocation() *LocationReport {
	t.updateIndex()

	now := t.getAdjustedTime()
	current := t.locationData.Destinations[t.currentIndex]
	next := t.next()

	log.Printf("Current: %s, Next: %s", current.City, next.City)

	report := new(LocationReport)
	report.Position = current.Location
	report.Next = next
	report.LastSeen = current
	report.CurrentTimeMs = now
	report.PresentsDelivered = t.calculatePresentsDelivered(now, current, next)

	// Check if we've departed (flying)
	if now > current.Departure && now < next.Arrival {
		report.Status = FLYING
		return report
	}

	// If we've not departed, we're at a stopover (landed)
	report.Status = LANDED
	return report
}

func NewTracker(xmasNow bool) (*Tracker, error) {
	var locationData LocationData

	// Read in location data from disk
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	file, err := os.Open(fmt.Sprintf("%s/data/location.json", path))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// FIXME: Update this to now use a non-deprecated method
	byteResult, _ := io.ReadAll(file)

	// Map the on disk locations to a JSON object
	err = json.Unmarshal(byteResult, &locationData)
	if err != nil {
		return nil, err
	}

	// Adjust the data start time to Christmas this year
	// TODO: Update this to overflow to next year if Santa's journey has ended?
	now := time.Now()
	dataStartTimeMs := locationData.Destinations[0].Arrival
	takeOffTimeThisYearMs := time.Date(now.Year(), time.December, 24, 10, 0, 0, 0, time.UTC).UnixMilli()
	timeDiff := takeOffTimeThisYearMs - dataStartTimeMs

	locationData.AdjustDepArrTimes(timeDiff)

	// If requested, make it Christmas now
	timeOffset := int64(0)
	if xmasNow && now.UnixMilli() < takeOffTimeThisYearMs {
		timeOffset = takeOffTimeThisYearMs - now.UnixMilli()
		//+ (time.Hour * 24).Milliseconds() + (time.Minute * 59).Milliseconds()
	}

	return &Tracker{0, &locationData, timeOffset}, nil
}
