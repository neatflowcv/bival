package entrygroup_test

import (
	"testing"
	"time"

	"github.com/neatflowcv/bival/internal/pkg/domain"
	"github.com/neatflowcv/bival/internal/pkg/domain/entrygroup"
	"github.com/stretchr/testify/require"
)

func issueCodes(issues []*entrygroup.Issue) []string {
	codes := make([]string, 0, len(issues))
	for _, issue := range issues {
		if issue == nil {
			continue
		}

		code := issue.Code()
		if code == "" {
			continue
		}

		codes = append(codes, code)
	}

	if len(codes) == 0 {
		return nil
	}

	return codes
}

func TestEntryGroupAddPlainTracksPendingMap(t *testing.T) {
	t.Parallel()

	group := entrygroup.New("alpha")

	group.AddPlain(newPlainEntry("alpha", true))
	require.True(t, group.HasPendingEntries())
}

func TestEntryGroupAddInstanceTracksPendingMap(t *testing.T) {
	t.Parallel()

	group := entrygroup.New("alpha")

	group.AddInstance(newInstanceEntry("alpha", true))
	require.True(t, group.HasPendingEntries())
}

func TestEntryGroupAddOLHTracksPendingLog(t *testing.T) {
	t.Parallel()

	group := entrygroup.New("alpha")

	group.AddOLH(newOLHEntry("alpha", true))
	require.True(t, group.HasPendingEntries())
}

func TestEntryGroupPendingStateMatchesStoredEntries(t *testing.T) {
	t.Parallel()

	group := entrygroup.New("alpha")

	require.False(t, group.HasPendingMap())
	require.False(t, group.HasPendingLog())
	require.False(t, group.HasPendingEntries())

	group.AddPlain(newPlainEntry("alpha", false))
	require.False(t, group.HasPendingMap())
	require.False(t, group.HasPendingLog())
	require.False(t, group.HasPendingEntries())

	group.AddInstance(newInstanceEntry("alpha", true))
	require.True(t, group.HasPendingMap())
	require.False(t, group.HasPendingLog())
	require.True(t, group.HasPendingEntries())

	otherGroup := entrygroup.New("beta")
	otherGroup.AddOLH(newOLHEntry("beta", true))
	require.False(t, otherGroup.HasPendingMap())
	require.True(t, otherGroup.HasPendingLog())
	require.True(t, otherGroup.HasPendingEntries())
}

func TestEntryGroupCountsMatchStoredEntries(t *testing.T) {
	t.Parallel()

	group := entrygroup.New("alpha")

	require.Zero(t, group.PlainCount())
	require.Zero(t, group.InstanceCount())
	require.Zero(t, group.OLHCount())

	group.AddPlain(newPlainEntry("alpha", false))
	require.Equal(t, len(group.PlainEntries()), group.PlainCount())

	group.AddPlain(newVersionedHeadPlainEntry())
	require.Equal(t, len(group.PlainEntries()), group.PlainCount())

	group.AddInstance(newInstanceEntry("alpha", false))
	require.Equal(t, len(group.InstanceEntries()), group.InstanceCount())

	group.AddOLH(newOLHEntry("alpha", false))
	require.Equal(t, len(group.OLHEntries()), group.OLHCount())
}

func TestEntryGroupPlainEntriesReturnsDefensiveCopy(t *testing.T) {
	t.Parallel()

	group := entrygroup.New("alpha")
	entry := newPlainEntry("alpha", false)
	group.AddPlain(entry)

	entries := group.PlainEntries()
	entries[0] = nil

	require.Len(t, group.PlainEntries(), 1)
	require.Same(t, entry, group.PlainEntries()[0])
}

func TestEntryGroupInstanceEntriesReturnsDefensiveCopy(t *testing.T) {
	t.Parallel()

	group := entrygroup.New("alpha")
	entry := newInstanceEntry("alpha", false)
	group.AddInstance(entry)

	entries := group.InstanceEntries()
	entries[0] = nil

	require.Len(t, group.InstanceEntries(), 1)
	require.Same(t, entry, group.InstanceEntries()[0])
}

func TestEntryGroupOLHEntriesReturnsDefensiveCopy(t *testing.T) {
	t.Parallel()

	group := entrygroup.New("alpha")
	entry := newOLHEntry("alpha", false)
	group.AddOLH(entry)

	entries := group.OLHEntries()
	entries[0] = nil

	require.Len(t, group.OLHEntries(), 1)
	require.Same(t, entry, group.OLHEntries()[0])
}

func TestEntryGroupProblemReasonReturnsEmptyForVersionedObjectAtThreshold(t *testing.T) {
	t.Parallel()

	group := entrygroup.New("alpha")

	versionOne := defaultVersionedFixture("v1")
	versionTwo := defaultVersionedFixture("v2")

	group.AddPlain(newVersionedHeadPlainEntry())
	group.AddPlain(newVersionedPlainEntry(versionOne))
	group.AddPlain(newVersionedPlainEntry(versionTwo))
	group.AddInstance(newVersionedInstanceEntry(versionOne))
	group.AddInstance(newVersionedInstanceEntry(versionTwo))
	group.AddOLH(newVersionedOLHEntry("alpha", "v1", false))
	require.Empty(t, group.ProblemReason())
}

func TestEntryGroupProblemReasonReportsTooManyVersionedEntries(t *testing.T) {
	t.Parallel()

	group := entrygroup.New("alpha")

	versionOne := defaultVersionedFixture("v1")
	versionTwo := defaultVersionedFixture("v2")
	versionThree := defaultVersionedFixture("v3")
	versionFour := defaultVersionedFixture("v4")

	group.AddPlain(newVersionedHeadPlainEntry())
	group.AddPlain(newVersionedPlainEntry(versionOne))
	group.AddPlain(newVersionedPlainEntry(versionTwo))
	group.AddPlain(newVersionedPlainEntry(versionThree))
	group.AddPlain(newVersionedPlainEntry(versionFour))
	group.AddInstance(newVersionedInstanceEntry(versionOne))
	group.AddInstance(newVersionedInstanceEntry(versionTwo))
	group.AddInstance(newVersionedInstanceEntry(versionThree))
	group.AddInstance(newVersionedInstanceEntry(versionFour))
	group.AddOLH(newVersionedOLHEntry("alpha", "v1", false))
	require.Equal(t, []string{"entry.versioned.count.exceeded"}, issueCodes(group.ProblemReason()))
}

func TestEntryGroupProblemReasonIncludesPendingEntry(t *testing.T) {
	t.Parallel()

	group := entrygroup.New("alpha")

	versionOne := fixtureWithPendingMap(defaultVersionedFixture("v1"), true)
	versionTwo := defaultVersionedFixture("v2")
	versionThree := defaultVersionedFixture("v3")

	group.AddPlain(newVersionedHeadPlainEntry())
	group.AddPlain(newVersionedPlainEntry(versionOne))
	group.AddPlain(newVersionedPlainEntry(versionTwo))
	group.AddPlain(newVersionedPlainEntry(versionThree))
	group.AddInstance(newVersionedInstanceEntry(versionOne))
	group.AddInstance(newVersionedInstanceEntry(versionTwo))
	group.AddInstance(newVersionedInstanceEntry(versionThree))
	group.AddOLH(newVersionedOLHEntry("alpha", "v1", false))
	require.Equal(t, []string{"plain.pending.exists", "instance.pending.exists"}, issueCodes(group.ProblemReason()))
}

func TestEntryGroupProblemReasonKeepsPendingEntryBeforeVersionedCountExceeded(t *testing.T) {
	t.Parallel()

	group := entrygroup.New("alpha")

	versionOne := fixtureWithPendingMap(defaultVersionedFixture("v1"), true)
	versionTwo := defaultVersionedFixture("v2")
	versionThree := defaultVersionedFixture("v3")
	versionFour := defaultVersionedFixture("v4")

	group.AddPlain(newVersionedHeadPlainEntry())
	group.AddPlain(newVersionedPlainEntry(versionOne))
	group.AddPlain(newVersionedPlainEntry(versionTwo))
	group.AddPlain(newVersionedPlainEntry(versionThree))
	group.AddPlain(newVersionedPlainEntry(versionFour))
	group.AddInstance(newVersionedInstanceEntry(versionOne))
	group.AddInstance(newVersionedInstanceEntry(versionTwo))
	group.AddInstance(newVersionedInstanceEntry(versionThree))
	group.AddInstance(newVersionedInstanceEntry(versionFour))
	group.AddOLH(newVersionedOLHEntry("alpha", "v1", false))

	require.Equal(
		t,
		[]string{"plain.pending.exists", "instance.pending.exists", "entry.versioned.count.exceeded"},
		issueCodes(group.ProblemReason()),
	)
}

func TestEntryGroupProblemReasonCountsPendingEntriesByField(t *testing.T) {
	t.Parallel()

	group := entrygroup.New("alpha")

	versionOne := fixtureWithPendingMap(defaultVersionedFixture("v1"), true)
	versionTwo := fixtureWithPendingMap(defaultVersionedFixture("v2"), true)
	versionThree := defaultVersionedFixture("v3")

	group.AddPlain(newVersionedHeadPlainEntry())
	group.AddPlain(newVersionedPlainEntry(versionOne))
	group.AddPlain(newVersionedPlainEntry(versionTwo))
	group.AddPlain(newVersionedPlainEntry(versionThree))
	group.AddInstance(newVersionedInstanceEntry(versionOne))
	group.AddInstance(newVersionedInstanceEntry(versionTwo))
	group.AddInstance(newVersionedInstanceEntry(versionThree))
	group.AddOLH(newVersionedOLHEntry("alpha", "v1", true))

	require.Equal(
		t,
		[]string{
			"plain.pending.exists",
			"plain.pending.exists",
			"instance.pending.exists",
			"instance.pending.exists",
			"olh.pending.exists",
			"olh.invalid",
		},
		issueCodes(group.ProblemReason()),
	)
}

func TestEntryGroupProblemReasonRejectsMultipleVersionsWhenOLHReferenceIsStale(t *testing.T) {
	t.Parallel()

	group := entrygroup.New("alpha")

	versionOne := defaultVersionedFixture("v1")
	versionOne.mtime = sampleStaleMTime()

	versionTwo := defaultVersionedFixture("v2")

	group.AddPlain(newVersionedHeadPlainEntry())
	group.AddPlain(newVersionedPlainEntry(versionOne))
	group.AddPlain(newVersionedPlainEntry(versionTwo))
	group.AddInstance(newVersionedInstanceEntry(versionOne))
	group.AddInstance(newVersionedInstanceEntry(versionTwo))
	group.AddOLH(newVersionedOLHEntry("alpha", "v1", false))
	require.Equal(t, []string{"olh.reference.outdated", "version.stale"}, issueCodes(group.ProblemReason()))
}

func TestEntryGroupProblemReasonAllowsOLHReferenceWithLatestMTime(t *testing.T) {
	t.Parallel()

	group := entrygroup.New("alpha")

	versionOne := defaultVersionedFixture("v1")
	versionOne.mtime = olderMTime

	versionTwo := defaultVersionedFixture("v2")
	versionTwo.mtime = newerMTime

	group.AddPlain(newVersionedHeadPlainEntry())
	group.AddPlain(newVersionedPlainEntry(versionOne))
	group.AddPlain(newVersionedPlainEntry(versionTwo))
	group.AddInstance(newVersionedInstanceEntry(versionOne))
	group.AddInstance(newVersionedInstanceEntry(versionTwo))
	group.AddOLH(newVersionedOLHEntry("alpha", "v2", false))

	require.Empty(t, group.ProblemReason())
}

func TestEntryGroupProblemReasonAllowsOLHReferenceWhenLatestMTimeIsTied(t *testing.T) {
	t.Parallel()

	group := entrygroup.New("alpha")

	versionOne := defaultVersionedFixture("v1")
	versionOne.mtime = newerMTime

	versionTwo := defaultVersionedFixture("v2")
	versionTwo.mtime = newerMTime

	group.AddPlain(newVersionedHeadPlainEntry())
	group.AddPlain(newVersionedPlainEntry(versionOne))
	group.AddPlain(newVersionedPlainEntry(versionTwo))
	group.AddInstance(newVersionedInstanceEntry(versionOne))
	group.AddInstance(newVersionedInstanceEntry(versionTwo))
	group.AddOLH(newVersionedOLHEntry("alpha", "v1", false))

	require.Empty(t, group.ProblemReason())
}

func TestEntryGroupProblemReasonSkipsOutdatedOLHReferenceWhenReferencedPairMTimeIsInvalid(t *testing.T) {
	t.Parallel()

	group := entrygroup.New("alpha")

	versionOne := defaultVersionedFixture("v1")
	versionOne.mtime = "invalid-mtime"

	versionTwo := defaultVersionedFixture("v2")
	versionTwo.mtime = newerMTime

	group.AddPlain(newVersionedHeadPlainEntry())
	group.AddPlain(newVersionedPlainEntry(versionOne))
	group.AddPlain(newVersionedPlainEntry(versionTwo))
	group.AddInstance(newVersionedInstanceEntry(versionOne))
	group.AddInstance(newVersionedInstanceEntry(versionTwo))
	group.AddOLH(newVersionedOLHEntry("alpha", "v1", false))

	require.Empty(t, group.ProblemReason())
}

func TestEntryGroupProblemReasonAllowsMultipleVersionsWhenOnlyNonOLHReferenceIsStale(t *testing.T) {
	t.Parallel()

	group := entrygroup.New("alpha")

	versionOne := defaultVersionedFixture("v1")
	versionOne.mtime = sampleStaleMTime()

	versionTwo := defaultVersionedFixture("v2")

	group.AddPlain(newVersionedHeadPlainEntry())
	group.AddPlain(newVersionedPlainEntry(versionOne))
	group.AddPlain(newVersionedPlainEntry(versionTwo))
	group.AddInstance(newVersionedInstanceEntry(versionOne))
	group.AddInstance(newVersionedInstanceEntry(versionTwo))
	group.AddOLH(newVersionedOLHEntry("alpha", "v2", false))
	require.Empty(t, group.ProblemReason())
}

func TestEntryGroupProblemReasonRejectsStaleDeleteMarkerOLHWhenVersionsExist(t *testing.T) {
	t.Parallel()

	group := entrygroup.New("alpha")

	version := defaultVersionedFixture("delete-v1")
	version.mtime = sampleStaleMTime()
	version.exists = false
	version.pool = -1
	version.epoch = 0
	version.eTag = ""
	version.tag = "delete-marker"
	version.flags = 7
	version.versionedEpoch = 3
	version.category = 0
	version.size = 0
	version.accountedSize = 0
	version.contentType = ""

	group.AddPlain(newVersionedHeadPlainEntry())
	group.AddPlain(newVersionedPlainEntry(version))
	group.AddInstance(newVersionedInstanceEntry(version))
	group.AddOLH(newCustomVersionedOLHEntry("alpha", "delete-v1", false, true))
	require.Equal(t, []string{"version.stale"}, issueCodes(group.ProblemReason()))
}

func TestEntryGroupRejectsMismatchedPlainName(t *testing.T) {
	t.Parallel()

	// Arrange
	group := entrygroup.New("alpha")

	// Act Assert
	require.PanicsWithValue(
		t,
		"entry name \"beta\" does not match group name \"alpha\"",
		func() {
			group.AddPlain(newPlainEntry("beta", false))
		},
	)
}

func TestEntryGroupRejectsMismatchedInstanceName(t *testing.T) {
	t.Parallel()

	// Arrange
	group := entrygroup.New("alpha")

	// Act Assert
	require.PanicsWithValue(
		t,
		"entry name \"beta\" does not match group name \"alpha\"",
		func() {
			group.AddInstance(newInstanceEntry("beta", false))
		},
	)
}

func TestEntryGroupRejectsMismatchedOLHName(t *testing.T) {
	t.Parallel()

	// Arrange
	group := entrygroup.New("alpha")

	// Act Assert
	require.PanicsWithValue(
		t,
		"entry name \"beta\" does not match group name \"alpha\"",
		func() {
			group.AddOLH(newOLHEntry("beta", false))
		},
	)
}

func newPlainEntry(name string, pending bool) *domain.Plain {
	return newCustomPlainEntry(
		name,
		name,
		"",
		pending,
		1,
		1,
		true,
		sampleMTime(),
		"etag",
		"tag",
		0,
	)
}

func newInstanceEntry(name string, pending bool) *domain.Instance {
	fixture := defaultVersionedFixture("v1")
	fixture.name = name

	return newVersionedInstanceEntry(fixtureWithPendingMap(fixture, pending))
}

type versionedEntryFixture struct {
	idx            string
	name           string
	instance       string
	pool           int
	epoch          int
	exists         bool
	mtime          string
	eTag           string
	tag            string
	flags          int
	versionedEpoch int
	category       int
	size           int64
	accountedSize  int64
	contentType    string
	owner          string
	ownerDisplay   string
	pendingMap     bool
}

func defaultVersionedFixture(instance string) versionedEntryFixture {
	return versionedEntryFixture{
		idx:            "alpha\x00v913\x00i" + instance,
		name:           "alpha",
		instance:       instance,
		pool:           186,
		epoch:          1147,
		exists:         true,
		mtime:          sampleMTime(),
		eTag:           "etag",
		tag:            "tag",
		flags:          1,
		versionedEpoch: 2,
		category:       1,
		size:           4,
		accountedSize:  4,
		contentType:    "text/plain",
		owner:          "test",
		ownerDisplay:   "test",
		pendingMap:     false,
	}
}

func fixtureWithPendingMap(fixture versionedEntryFixture, pending bool) versionedEntryFixture {
	fixture.pendingMap = pending

	return fixture
}

func newVersionedHeadPlainEntry() *domain.Plain {
	return newCustomVersionedPlainEntry(versionedEntryFixture{
		idx:            "alpha",
		name:           "alpha",
		instance:       "",
		pool:           -1,
		epoch:          0,
		exists:         false,
		mtime:          "0.000000",
		eTag:           "",
		tag:            "",
		flags:          8,
		versionedEpoch: 0,
		category:       0,
		size:           0,
		accountedSize:  0,
		contentType:    "",
		owner:          "",
		ownerDisplay:   "",
		pendingMap:     false,
	}, false)
}

func newVersionedPlainEntry(fixture versionedEntryFixture) *domain.Plain {
	return newCustomVersionedPlainEntry(fixture, true)
}

func newCustomVersionedPlainEntry(fixture versionedEntryFixture, buildVersionedIndex bool) *domain.Plain {
	var pendingMaps []*domain.PendingMap
	if fixture.pendingMap {
		pendingMaps = []*domain.PendingMap{nil}
	}

	idx := fixture.idx
	if buildVersionedIndex {
		idx = versionedPlainIndex(fixture.instance)
	}

	return domain.NewPlain(
		domain.DirEntryParams{
			Kind:             "plain",
			Index:            []byte(idx),
			Name:             fixture.name,
			Instance:         fixture.instance,
			Pool:             fixture.pool,
			Epoch:            fixture.epoch,
			VEpoch:           fixture.versionedEpoch,
			Locator:          "",
			Exists:           fixture.exists,
			Tag:              fixture.tag,
			Flags:            fixture.flags,
			Category:         fixture.category,
			Size:             fixture.size,
			AccountedSize:    fixture.accountedSize,
			Appendable:       false,
			MTime:            fixture.mtime,
			ETag:             fixture.eTag,
			StorageClass:     "",
			ContentType:      fixture.contentType,
			OwnerUserID:      fixture.owner,
			OwnerDisplayName: fixture.ownerDisplay,
			PendingMaps:      pendingMaps,
		},
	)
}

func newVersionedInstanceEntry(fixture versionedEntryFixture) *domain.Instance {
	var pendingMaps []*domain.PendingMap
	if fixture.pendingMap {
		pendingMaps = []*domain.PendingMap{nil}
	}

	return domain.NewInstance(
		domain.DirEntryParams{
			Kind:             "instance",
			Index:            []byte(versionedInstanceIndex(fixture.instance)),
			Name:             fixture.name,
			Instance:         fixture.instance,
			Pool:             fixture.pool,
			Epoch:            fixture.epoch,
			VEpoch:           fixture.versionedEpoch,
			Locator:          "",
			Exists:           fixture.exists,
			Tag:              fixture.tag,
			Flags:            fixture.flags,
			Category:         fixture.category,
			Size:             fixture.size,
			AccountedSize:    fixture.accountedSize,
			Appendable:       false,
			MTime:            fixture.mtime,
			ETag:             fixture.eTag,
			StorageClass:     "",
			ContentType:      fixture.contentType,
			OwnerUserID:      fixture.owner,
			OwnerDisplayName: fixture.ownerDisplay,
			PendingMaps:      pendingMaps,
		},
	)
}

func newCustomPlainEntry(
	idx string,
	name string,
	instance string,
	pending bool,
	pool int,
	epoch int,
	exists bool,
	mtime string,
	etag string,
	tag string,
	flags int,
) *domain.Plain {
	var pendingMaps []*domain.PendingMap
	if pending {
		pendingMaps = []*domain.PendingMap{nil}
	}

	return domain.NewPlain(
		domain.DirEntryParams{
			Kind:             "plain",
			Index:            []byte(idx),
			Name:             name,
			Instance:         instance,
			Pool:             pool,
			Epoch:            epoch,
			VEpoch:           0,
			Locator:          "",
			Exists:           exists,
			Tag:              tag,
			Flags:            flags,
			Category:         1,
			Size:             4,
			AccountedSize:    4,
			Appendable:       false,
			MTime:            mtime,
			ETag:             etag,
			StorageClass:     "",
			ContentType:      "",
			OwnerUserID:      "",
			OwnerDisplayName: "",
			PendingMaps:      pendingMaps,
		},
	)
}

func sampleMTime() string {
	return sampleMTimeAtOffset(-24 * time.Hour)
}

func sampleStaleMTime() string {
	return sampleMTimeAtOffset(-8 * 24 * time.Hour)
}

func sampleMTimeAtOffset(offset time.Duration) string {
	return time.Now().UTC().Truncate(24 * time.Hour).Add(offset).Format(time.RFC3339Nano)
}

func newOLHEntry(name string, pending bool) *domain.OLH {
	return newVersionedOLHEntry(name, "v1", pending)
}

func newVersionedOLHEntry(name string, instance string, pending bool) *domain.OLH {
	return newCustomVersionedOLHEntry(name, instance, pending, false)
}

func newCustomVersionedOLHEntry(name string, instance string, pending bool, deleteMarker bool) *domain.OLH {
	var pendingLogs []*domain.PendingLog
	if pending {
		pendingLogs = []*domain.PendingLog{nil}
	}

	return domain.NewOLH(domain.OLHParams{
		Kind:           "olh",
		Index:          []byte(name),
		Name:           name,
		Instance:       instance,
		DeleteMarker:   deleteMarker,
		PendingRemoval: false,
		Exists:         false,
		Epoch:          0,
		PendingLogs:    pendingLogs,
		Tag:            "",
	})
}

func versionedPlainIndex(instance string) string {
	return "alpha\x00v913\x00i" + instance
}

func versionedInstanceIndex(instance string) string {
	return "\x801000_alpha\x00i" + instance
}
