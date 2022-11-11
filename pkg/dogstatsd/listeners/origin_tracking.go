package listeners

import (
	"github.com/DataDog/datadog-agent/pkg/telemetry"
)

var (
	tlmOriginTrackingBytes = telemetry.NewCounter("dogstatsd", "origin_tracking_bytes",
		[]string{"origin"}, "Bytes per origin")
	tlmOriginTrackingPackets = telemetry.NewCounter("dogstatsd", "origin_tracking_packets",
		[]string{"origin"}, "Packets per origin")
	tlmOriginTrackingTags = telemetry.NewCounter("dogstatsd", "origin_tracking_tags",
		[]string{"origin"}, "Tags per origin")
)

type OriginTrackingEntry struct {
	Origin     string
	BytesCount uint
	TagsCount  uint
}

type OriginTracker struct {
	ch           chan OriginTrackingEntry
	packetsCount map[string]uint
	bytesCount   map[string]uint
	tagsCount    map[string]uint
}

func StartOriginTracker() *OriginTracker {
	trackingCh := make(chan OriginTrackingEntry, 8192)
	tracker := &OriginTracker{
		ch:           trackingCh,
		packetsCount: make(map[string]uint),
		bytesCount:   make(map[string]uint),
		tagsCount:    make(map[string]uint),
	}
	go func() {
		tracker.runOriginTracking(trackingCh)
	}()
	return tracker
}

func (t *OriginTracker) processEntry(entry OriginTrackingEntry) {
	tlmOriginTrackingBytes.Add(float64(entry.BytesCount), entry.Origin)
	tlmOriginTrackingTags.Add(float64(entry.TagsCount), entry.Origin)
	tlmOriginTrackingPackets.Inc(entry.Origin)
}

func countCommas(b []byte) uint {
	var rv uint = 0
	for i := 0; i < len(b); i++ {
		if b[i] == ',' {
			rv += 1
		}
	}
	return rv
}

func (t *OriginTracker) runOriginTracking(trackingCh chan OriginTrackingEntry) {
	for {
		select {
		case entry := <-trackingCh:
			t.processEntry(entry)
		}
	}
}

// TODO(remy): close channel and stop the tracker
