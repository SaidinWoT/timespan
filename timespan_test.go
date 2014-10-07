package timespan_test

import (
	"testing"
	"time"
	"timespan"
)

var times = []time.Time{
	time.Date(2014, time.February, 3, 2, 0, 0, 0, time.UTC),
	time.Date(2014, time.February, 3, 4, 0, 0, 0, time.UTC),
	time.Date(2014, time.February, 3, 6, 0, 0, 0, time.UTC),
	time.Date(2014, time.February, 3, 8, 0, 0, 0, time.UTC),
}

var durations = []time.Duration{
	time.Duration(2) * time.Hour,
	time.Duration(4) * time.Hour,
	time.Duration(6) * time.Hour,
}

var spans = []timespan.Span{
	timespan.NewSpan(times[0], durations[0]), // 2:00 - 4:00
	timespan.NewSpan(times[0], durations[1]), // 2:00 - 6:00
	timespan.NewSpan(times[0], durations[2]), // 2:00 - 8:00
	timespan.NewSpan(times[1], durations[0]), // 4:00 - 6:00
	timespan.NewSpan(times[1], durations[1]), // 4:00 - 8:00
	timespan.NewSpan(times[2], durations[0]), // 6:00 - 8:00
}

func TestNewSpan(t *testing.T) {
	if spans[0].Start() != times[0] {
		t.Error("Improper timespan start value.")
	}
	if spans[0].End() != times[1] {
		t.Error("Improper timespan end value.")
	}
}

func TestAfter(t *testing.T) {
	if spans[5].After(times[3]) { // 6:00 - 8:00 >? 8:00
		t.Error("Span reported as after its end time.")
	}
	if spans[5].After(times[2]) { // 6:00 - 8:00 >? 6:00
		t.Error("Span reported as after its start time.")
	}
	if !spans[5].After(times[1]) { // 6:00 - 8:00 >? 4:00
		t.Error("Span reported as not after earlier time.")
	}
}

func TestBefore(t *testing.T) {
	if spans[0].Before(times[0]) { // 2:00 - 4:00 <? 2:00
		t.Error("Span reported as before its start time.")
	}
	if spans[0].Before(times[1]) { // 2:00 - 4:00 <? 4:00
		t.Error("Span reported as before its end time.")
	}
	if !spans[0].Before(times[2]) { // 2:00 - 4:00 <? 6:00
		t.Error("Span reported as not before later time.")
	}
}

func TestFollows(t *testing.T) {
	if !spans[3].Follows(spans[0]) { // 4:00 - 6:00 >? 2:00 - 4:00
		t.Error("Span reported as not following an earlier span.")
	}
	if spans[3].Follows(spans[1]) { // 4:00 - 6:00 >? 2:00 - 6:00
		t.Error("Span reported as following a containing span ending at the same time.")
	}
	if spans[3].Follows(spans[3]) { // 4:00 - 6:00 >? 4:00 - 6:00
		t.Error("Span reported as following itself.")
	}
	if spans[3].Follows(spans[4]) { // 4:00 - 6:00 >? 4:00 - 8:00
		t.Error("Span reported as following a containing span ending at a later time.")
	}
	if spans[3].Follows(spans[5]) { // 4:00 - 6:00 >? 6:00 - 8:00
		t.Error("Span reported as following a later span.")
	}
	if spans[2].Follows(spans[3]) { // 2:00 - 8:00 >? 4:00 - 6:00
		t.Error("Span reported as following a fully contained span.")
	}
	if spans[4].Follows(spans[1]) { // 4:00 - 8:00 >? 2:00 - 6:00
		t.Error("Span reported as following an overlapping span ending earlier.")
	}
	if spans[1].Follows(spans[4]) { // 2:00 - 6:00 >? 4:00 - 8:00
		t.Error("Span reported as following an overlapping span ending later.")
	}
}

func TestPrecedes(t *testing.T) {
	if spans[3].Precedes(spans[0]) { // 4:00 - 6:00 <? 2:00 - 4:00
		t.Error("Span reported as preceding an earlier span.")
	}
	if spans[3].Precedes(spans[1]) { // 4:00 - 6:00 <? 2:00 - 6:00
		t.Error("Span reported as preceding a containing span starting earlier.")
	}
	if spans[3].Precedes(spans[3]) { // 4:00 - 6:00 <? 4:00 - 6:00
		t.Error("Span reported as preceding itself.")
	}
	if spans[3].Precedes(spans[4]) { // 4:00 - 6:00 <? 4:00 - 8:00
		t.Error("Span reported as preceding a containing span starting at the same time.")
	}
	if !spans[3].Precedes(spans[5]) { // 4:00 - 6:00 <? 6:00 - 8:00
		t.Error("Span reported as not preceding a later span.")
	}
	if spans[2].Precedes(spans[3]) { // 2:00 - 8:00 <? 4:00 - 6:00
		t.Error("Span reported as preceding a fully contained span.")
	}
	if spans[4].Precedes(spans[1]) { // 4:00 - 8:00 <? 2:00 - 6:00
		t.Error("Span reported as preceding an overlapping span starting earlier.")
	}
	if spans[1].Precedes(spans[4]) { // 2:00 - 6:00 <? 4:00 - 8:00
		t.Error("Span reported as preceding an overlapping span starting later.")
	}
}

func TestEqual(t *testing.T) {
	if !spans[3].Equal(spans[3]) { // 4:00 - 6:00 =? 4:00 - 6:00
		t.Error("Span reported as not equal to itself.")
	}
	if spans[3].Equal(spans[2]) { // 4:00 - 6:00 =? 2:00 - 8:00
		t.Error("Span reported as equal to a fully containing span.")
	}
	if spans[2].Equal(spans[3]) { // 2:00 - 8:00 =? 4:00 - 6:00
		t.Error("Span reported as equal to a fully contained span.")
	}
	if spans[3].Equal(spans[1]) { // 4:00 - 6:00 <? 2:00 - 6:00
		t.Error("Span reported as equal to span with different start time.")
	}
	if spans[3].Equal(spans[4]) { // 4:00 - 6:00 <? 4:00 - 8:00
		t.Error("Span reported as equal to span with different end time.")
	}
	if spans[1].Equal(spans[4]) { // 2:00 - 6:00 =? 4:00 - 8:00
		t.Error("Span reported as equal to overlapping span.")
	}
}

var (
    d time.Duration
    r timespan.Span
    s timespan.Span
    t time.Time
)

func BenchmarkStart(b *testing.B) {
    s = spans[0]
    for i := 0; i < b.N; i++ {
        t = s.Start()
    }
}

func BenchmarkEnd(b *testing.B) {
    s = spans[0]
    for i := 0; i < b.N; i++ {
        t = s.End()
    }
}

func BenchmarkDuration(b *testing.B) {
    s = spans[0]
    for i := 0; i < b.N; i++ {
        d = s.Duration()
    }
}

func BenchmarkAfter(b *testing.B) {
    s = spans[5]
    t = times[0]
    for i := 0; i < b.N; i++ {
        _ = s.After(t)
    }
}

func BenchmarkBefore(b *testing.B) {
    s = spans[5]
    t = times[0]
    for i := 0; i < b.N; i++ {
        _ = s.Before(t)
    }
}

func BenchmarkFollows(b *testing.B) {
    s, r = spans[5], spans[0]
    for i := 0; i < b.N; i++ {
        _ = s.Follows(r)
    }
}

func BenchmarkPrecedes(b *testing.B) {
    s, r = spans[5], spans[0]
    for i := 0; i < b.N; i++ {
        _ = s.Precedes(r)
    }
}

func BenchmarkContainsTime(b *testing.B) {
    s = spans[2]
    t = times[1]
    for i := 0; i < b.N; i++ {
        _ = s.ContainsTime(t)
    }
}

func BenchmarkContains(b *testing.B) {
    s, r = spans[2], spans[3]
    for i := 0; i < b.N; i++ {
        _ = s.Contains(r)
    }
}

func BenchmarkEncompass(b *testing.B) {
    s, r = spans[0], spans[5]
    for i := 0; i < b.N; i++ {
        _ = s.Encompass(r)
    }
}

func BenchmarkGap(b *testing.B) {
    s, r = spans[0], spans[5]
    for i := 0; i < b.N; i++ {
        _ = s.Gap(r)
    }
}

func BenchmarkIntersection(b *testing.B) {
    s, r = spans[1], spans[4]
    for i := 0; i < b.N; i++ {
        _ = s.Intersection(r)
    }
}

func BenchmarkOffset(b *testing.B) {
    s = spans[0]
    d = durations[1]
    for i := 0; i < b.N; i++ {
        r = s.Offset(d)
    }
}

func BenchmarkOffsetDate(b *testing.B) {
    s = spans[0]
    for i := 0; i < b.N; i++ {
        r = s.OffsetDate(1, 1, 1)
    }
}

func BenchmarkOverlaps(b *testing.B) {
    s, r = spans[0], spans[1]
    for i := 0; i < b.N; i++ {
        _ = s.Overlaps(r)
    }
}

func BenchmarkIsZero(b *testing.B) {
    s = spans[0]
    for i := 0; i < b.N; i++ {
        _ = s.IsZero()
    }
}

func BenchmarkEqual(b *testing.B) {
    s, r = spans[0], spans[1]
    for i := 0; i < b.N; i++ {
        _ = s.Equal(r)
    }
}

func BenchmarkBorders(b *testing.B) {
    s, r = spans[0], spans[3]
    for i := 0; i < b.N; i++ {
        _ = s.Borders(r)
    }
}
