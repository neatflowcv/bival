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

	err := group.AddPlain(newPlainEntry("alpha", true))
	require.NoError(t, err)
	require.True(t, group.HasPendingEntries())
}

func TestEntryGroupAddInstanceTracksPendingMap(t *testing.T) {
	t.Parallel()

	group := entrygroup.New("alpha")

	err := group.AddInstance(newInstanceEntry("alpha", true))
	require.NoError(t, err)
	require.True(t, group.HasPendingEntries())
}

func TestEntryGroupAddOLHTracksPendingLog(t *testing.T) {
	t.Parallel()

	group := entrygroup.New("alpha")

	err := group.AddOLH(newOLHEntry("alpha", true))
	require.NoError(t, err)
	require.True(t, group.HasPendingEntries())
}

func TestEntryGroupPendingStateMatchesStoredEntries(t *testing.T) {
	t.Parallel()

	group := entrygroup.New("alpha")

	require.False(t, group.HasPendingMap())
	require.False(t, group.HasPendingLog())
	require.False(t, group.HasPendingEntries())

	require.NoError(t, group.AddPlain(newPlainEntry("alpha", false)))
	require.False(t, group.HasPendingMap())
	require.False(t, group.HasPendingLog())
	require.False(t, group.HasPendingEntries())

	require.NoError(t, group.AddInstance(newInstanceEntry("alpha", true)))
	require.True(t, group.HasPendingMap())
	require.False(t, group.HasPendingLog())
	require.True(t, group.HasPendingEntries())

	otherGroup := entrygroup.New("beta")
	require.NoError(t, otherGroup.AddOLH(newOLHEntry("beta", true)))
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

	require.NoError(t, group.AddPlain(newPlainEntry("alpha", false)))
	require.Equal(t, len(group.PlainEntries()), group.PlainCount())

	require.NoError(t, group.AddPlain(newVersionedHeadPlainEntry()))
	require.Equal(t, len(group.PlainEntries()), group.PlainCount())

	require.NoError(t, group.AddInstance(newInstanceEntry("alpha", false)))
	require.Equal(t, len(group.InstanceEntries()), group.InstanceCount())

	require.NoError(t, group.AddOLH(newOLHEntry("alpha", false)))
	require.Equal(t, len(group.OLHEntries()), group.OLHCount())
}

func TestEntryGroupProblemReasonReturnsEmptyForVersionedObjectAtThreshold(t *testing.T) {
	t.Parallel()

	group := entrygroup.New("alpha")

	versionOne := defaultVersionedFixture("v1")
	versionTwo := defaultVersionedFixture("v2")

	require.NoError(t, group.AddPlain(newVersionedHeadPlainEntry()))
	require.NoError(t, group.AddPlain(newVersionedPlainEntry(versionOne)))
	require.NoError(t, group.AddPlain(newVersionedPlainEntry(versionTwo)))
	require.NoError(t, group.AddInstance(newVersionedInstanceEntry(versionOne)))
	require.NoError(t, group.AddInstance(newVersionedInstanceEntry(versionTwo)))
	require.NoError(t, group.AddOLH(newVersionedOLHEntry("alpha", "v1", false)))
	require.Empty(t, group.ProblemReason())
}

func TestEntryGroupProblemReasonReportsTooManyVersionedEntries(t *testing.T) {
	t.Parallel()

	group := entrygroup.New("alpha")

	versionOne := defaultVersionedFixture("v1")
	versionTwo := defaultVersionedFixture("v2")
	versionThree := defaultVersionedFixture("v3")
	versionFour := defaultVersionedFixture("v4")

	require.NoError(t, group.AddPlain(newVersionedHeadPlainEntry()))
	require.NoError(t, group.AddPlain(newVersionedPlainEntry(versionOne)))
	require.NoError(t, group.AddPlain(newVersionedPlainEntry(versionTwo)))
	require.NoError(t, group.AddPlain(newVersionedPlainEntry(versionThree)))
	require.NoError(t, group.AddPlain(newVersionedPlainEntry(versionFour)))
	require.NoError(t, group.AddInstance(newVersionedInstanceEntry(versionOne)))
	require.NoError(t, group.AddInstance(newVersionedInstanceEntry(versionTwo)))
	require.NoError(t, group.AddInstance(newVersionedInstanceEntry(versionThree)))
	require.NoError(t, group.AddInstance(newVersionedInstanceEntry(versionFour)))
	require.NoError(t, group.AddOLH(newVersionedOLHEntry("alpha", "v1", false)))
	require.Equal(t, []string{"entry.versioned.count.exceeded"}, issueCodes(group.ProblemReason()))
}

func TestEntryGroupProblemReasonIncludesPendingEntry(t *testing.T) {
	t.Parallel()

	group := entrygroup.New("alpha")

	versionOne := fixtureWithPendingMap(defaultVersionedFixture("v1"), true)
	versionTwo := defaultVersionedFixture("v2")
	versionThree := defaultVersionedFixture("v3")

	require.NoError(t, group.AddPlain(newVersionedHeadPlainEntry()))
	require.NoError(t, group.AddPlain(newVersionedPlainEntry(versionOne)))
	require.NoError(t, group.AddPlain(newVersionedPlainEntry(versionTwo)))
	require.NoError(t, group.AddPlain(newVersionedPlainEntry(versionThree)))
	require.NoError(t, group.AddInstance(newVersionedInstanceEntry(versionOne)))
	require.NoError(t, group.AddInstance(newVersionedInstanceEntry(versionTwo)))
	require.NoError(t, group.AddInstance(newVersionedInstanceEntry(versionThree)))
	require.NoError(t, group.AddOLH(newVersionedOLHEntry("alpha", "v1", false)))
	require.Equal(t, []string{"entry.pending.exists"}, issueCodes(group.ProblemReason()))
}

func TestEntryGroupProblemReasonKeepsPendingEntryBeforeVersionedCountExceeded(t *testing.T) {
	t.Parallel()

	group := entrygroup.New("alpha")

	versionOne := fixtureWithPendingMap(defaultVersionedFixture("v1"), true)
	versionTwo := defaultVersionedFixture("v2")
	versionThree := defaultVersionedFixture("v3")
	versionFour := defaultVersionedFixture("v4")

	require.NoError(t, group.AddPlain(newVersionedHeadPlainEntry()))
	require.NoError(t, group.AddPlain(newVersionedPlainEntry(versionOne)))
	require.NoError(t, group.AddPlain(newVersionedPlainEntry(versionTwo)))
	require.NoError(t, group.AddPlain(newVersionedPlainEntry(versionThree)))
	require.NoError(t, group.AddPlain(newVersionedPlainEntry(versionFour)))
	require.NoError(t, group.AddInstance(newVersionedInstanceEntry(versionOne)))
	require.NoError(t, group.AddInstance(newVersionedInstanceEntry(versionTwo)))
	require.NoError(t, group.AddInstance(newVersionedInstanceEntry(versionThree)))
	require.NoError(t, group.AddInstance(newVersionedInstanceEntry(versionFour)))
	require.NoError(t, group.AddOLH(newVersionedOLHEntry("alpha", "v1", false)))

	require.Equal(
		t,
		[]string{"entry.pending.exists", "entry.versioned.count.exceeded"},
		issueCodes(group.ProblemReason()),
	)
}

func TestEntryGroupProblemReasonRejectsMultipleVersionsWhenOLHReferenceIsStale(t *testing.T) {
	t.Parallel()

	group := entrygroup.New("alpha")

	versionOne := defaultVersionedFixture("v1")
	versionOne.mtime = sampleStaleMTime()

	versionTwo := defaultVersionedFixture("v2")

	require.NoError(t, group.AddPlain(newVersionedHeadPlainEntry()))
	require.NoError(t, group.AddPlain(newVersionedPlainEntry(versionOne)))
	require.NoError(t, group.AddPlain(newVersionedPlainEntry(versionTwo)))
	require.NoError(t, group.AddInstance(newVersionedInstanceEntry(versionOne)))
	require.NoError(t, group.AddInstance(newVersionedInstanceEntry(versionTwo)))
	require.NoError(t, group.AddOLH(newVersionedOLHEntry("alpha", "v1", false)))
	require.Equal(t, []string{"olh.reference.stale"}, issueCodes(group.ProblemReason()))
}

func TestEntryGroupProblemReasonAllowsMultipleVersionsWhenOnlyNonOLHReferenceIsStale(t *testing.T) {
	t.Parallel()

	group := entrygroup.New("alpha")

	versionOne := defaultVersionedFixture("v1")
	versionOne.mtime = sampleStaleMTime()

	versionTwo := defaultVersionedFixture("v2")

	require.NoError(t, group.AddPlain(newVersionedHeadPlainEntry()))
	require.NoError(t, group.AddPlain(newVersionedPlainEntry(versionOne)))
	require.NoError(t, group.AddPlain(newVersionedPlainEntry(versionTwo)))
	require.NoError(t, group.AddInstance(newVersionedInstanceEntry(versionOne)))
	require.NoError(t, group.AddInstance(newVersionedInstanceEntry(versionTwo)))
	require.NoError(t, group.AddOLH(newVersionedOLHEntry("alpha", "v2", false)))
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

	require.NoError(t, group.AddPlain(newVersionedHeadPlainEntry()))
	require.NoError(t, group.AddPlain(newVersionedPlainEntry(version)))
	require.NoError(t, group.AddInstance(newVersionedInstanceEntry(version)))
	require.NoError(t, group.AddOLH(newCustomVersionedOLHEntry("alpha", "delete-v1", false, true)))
	require.Equal(t, []string{"olh.delete_marker.stale"}, issueCodes(group.ProblemReason()))
}

func TestEntryGroupClassifierReturnsVersionedObjectWhenNoRuleMatches(t *testing.T) {
	t.Parallel()

	group := entrygroup.New("alpha")

	require.NoError(t, group.AddPlain(newUnversionedPlainEntry()))
	require.NoError(t, group.AddPlain(newVersionedPlainEntry(defaultVersionedFixture("v1"))))
	require.NoError(t, group.AddInstance(newVersionedInstanceEntry(defaultVersionedFixture("v1"))))
}

func TestEntryGroupClassifierReturnsVersionedObjectWithSingleVersion(t *testing.T) {
	t.Parallel()

	group := entrygroup.New("alpha")

	version := defaultVersionedFixture("v1")
	instanceVersion := version
	instanceVersion.tag = "instance-tag"

	require.NoError(t, group.AddPlain(newVersionedHeadPlainEntry()))
	require.NoError(t, group.AddPlain(newVersionedPlainEntry(version)))
	require.NoError(t, group.AddInstance(newVersionedInstanceEntry(instanceVersion)))
	require.NoError(t, group.AddOLH(newVersionedOLHEntry("alpha", "v1", false)))
}

func TestEntryGroupClassifierReturnsVersionedObjectWithDeleteMarkerHead(t *testing.T) {
	t.Parallel()

	group := entrygroup.New("alpha")

	deleteMarker := defaultVersionedFixture("delete-v1")
	deleteMarker.exists = false
	deleteMarker.pool = -1
	deleteMarker.epoch = 0
	deleteMarker.eTag = ""
	deleteMarker.tag = "delete-marker"
	deleteMarker.flags = 7
	deleteMarker.versionedEpoch = 3
	deleteMarker.category = 0
	deleteMarker.size = 0
	deleteMarker.accountedSize = 0
	deleteMarker.contentType = ""
	deleteMarker.mtime = sampleDeleteMarkerMTime()

	require.NoError(t, group.AddPlain(newVersionedHeadPlainEntry()))
	require.NoError(t, group.AddPlain(newVersionedPlainEntry(deleteMarker)))
	require.NoError(t, group.AddInstance(newVersionedInstanceEntry(deleteMarker)))
	require.NoError(t, group.AddOLH(newVersionedOLHEntry("alpha", "delete-v1", false)))
}

func TestEntryGroupClassifierReturnsVersionedObjectWithMultipleVersions(t *testing.T) {
	t.Parallel()

	group := entrygroup.New("alpha")

	versionOne := defaultVersionedFixture("v1")
	versionOne.pool = 186
	versionOne.epoch = 1121
	versionOne.tag = "tag-v1"
	versionOne.flags = 3
	versionOne.versionedEpoch = 3
	versionOne.mtime = sampleEarlierMTime()

	versionTwo := defaultVersionedFixture("v2")

	require.NoError(t, group.AddPlain(newVersionedHeadPlainEntry()))
	require.NoError(t, group.AddPlain(newVersionedPlainEntry(versionOne)))
	require.NoError(t, group.AddPlain(newVersionedPlainEntry(versionTwo)))
	require.NoError(t, group.AddInstance(newVersionedInstanceEntry(versionOne)))
	require.NoError(t, group.AddInstance(newVersionedInstanceEntry(versionTwo)))
	require.NoError(t, group.AddOLH(newVersionedOLHEntry("alpha", "v1", false)))
}

func TestEntryGroupClassifierReturnsVersionedObjectWithEmptyInstanceVersion(t *testing.T) {
	t.Parallel()

	group := entrygroup.New("alpha")

	version := defaultVersionedFixture("")
	version.pool = 34
	version.epoch = 2219298
	version.versionedEpoch = 3
	version.flags = 3

	require.NoError(t, group.AddPlain(newVersionedHeadPlainEntry()))
	require.NoError(t, group.AddPlain(newVersionedPlainEntry(version)))
	require.NoError(t, group.AddInstance(newVersionedInstanceEntry(version)))
	require.NoError(t, group.AddOLH(newVersionedOLHEntry("alpha", "", false)))
}

func TestEntryGroupClassifierRejectsVersionedObjectWhenHeadPlainIsMissing(t *testing.T) {
	t.Parallel()

	group := entrygroup.New("alpha")

	require.NoError(t, group.AddPlain(newVersionedPlainEntry(defaultVersionedFixture("v1"))))
	require.NoError(t, group.AddInstance(newVersionedInstanceEntry(defaultVersionedFixture("v1"))))
	require.NoError(t, group.AddOLH(newVersionedOLHEntry("alpha", "v1", false)))
}

func TestEntryGroupClassifierRejectsVersionedObjectWhenHeadPlainIsDuplicated(t *testing.T) {
	t.Parallel()

	group := entrygroup.New("alpha")

	require.NoError(t, group.AddPlain(newVersionedHeadPlainEntry()))
	require.NoError(t, group.AddPlain(newVersionedHeadPlainEntry()))
	require.NoError(t, group.AddInstance(newVersionedInstanceEntry(defaultVersionedFixture("v1"))))
	require.NoError(t, group.AddOLH(newVersionedOLHEntry("alpha", "v1", false)))
}

func TestEntryGroupClassifierRejectsVersionedObjectWhenHeadPlainShapeIsInvalid(t *testing.T) {
	t.Parallel()

	fixture := versionedHeadFixture()
	fixture.flags = 0

	assertVersionedObjectRejectedForInvalidHead(t, fixture)
}

func TestEntryGroupClassifierRejectsVersionedObjectWhenHeadPlainMetaIsNotDefault(t *testing.T) {
	t.Parallel()

	fixture := versionedHeadFixture()
	fixture.category = 1

	assertVersionedObjectRejectedForInvalidHead(t, fixture)
}

func assertVersionedObjectRejectedForInvalidHead(t *testing.T, fixture versionedEntryFixture) {
	t.Helper()

	group := entrygroup.New("alpha")

	invalidHead := newCustomVersionedPlainEntry(fixture, false)

	require.NoError(t, group.AddPlain(invalidHead))
	require.NoError(t, group.AddPlain(newVersionedPlainEntry(defaultVersionedFixture("v1"))))
	require.NoError(t, group.AddInstance(newVersionedInstanceEntry(defaultVersionedFixture("v1"))))
	require.NoError(t, group.AddOLH(newVersionedOLHEntry("alpha", "v1", false)))
}

func versionedHeadFixture() versionedEntryFixture {
	return versionedEntryFixture{
		idx:            "alpha",
		name:           "alpha",
		instance:       "",
		pool:           -1,
		epoch:          0,
		exists:         false,
		mtime:          "0.000000",
		eTag:           "",
		tag:            "",
		flags:          0,
		versionedEpoch: 0,
		category:       0,
		size:           0,
		accountedSize:  0,
		contentType:    "",
		owner:          "",
		ownerDisplay:   "",
		pendingMap:     false,
	}
}

func TestEntryGroupClassifierRejectsVersionedObjectWhenPlainInstanceNonTagPayloadsDiffer(t *testing.T) {
	t.Parallel()

	group := entrygroup.New("alpha")

	fixture := defaultVersionedFixture("v1")
	mismatchedInstance := fixture
	mismatchedInstance.flags = fixture.flags + 1

	require.NoError(t, group.AddPlain(newVersionedHeadPlainEntry()))
	require.NoError(t, group.AddPlain(newVersionedPlainEntry(fixture)))
	require.NoError(t, group.AddInstance(newVersionedInstanceEntry(mismatchedInstance)))
	require.NoError(t, group.AddOLH(newVersionedOLHEntry("alpha", "v1", false)))
}

func TestEntryGroupClassifierRejectsVersionedObjectWhenPlainHasNoMatchingInstance(t *testing.T) {
	t.Parallel()

	group := entrygroup.New("alpha")

	require.NoError(t, group.AddPlain(newVersionedHeadPlainEntry()))
	require.NoError(t, group.AddPlain(newVersionedPlainEntry(defaultVersionedFixture("v1"))))
	require.NoError(t, group.AddPlain(newVersionedPlainEntry(defaultVersionedFixture("v2"))))
	require.NoError(t, group.AddInstance(newVersionedInstanceEntry(defaultVersionedFixture("v1"))))
	require.NoError(t, group.AddOLH(newVersionedOLHEntry("alpha", "v1", false)))
}

func TestEntryGroupClassifierRejectsVersionedObjectWhenInstanceHasNoMatchingPlain(t *testing.T) {
	t.Parallel()

	group := entrygroup.New("alpha")

	require.NoError(t, group.AddPlain(newVersionedHeadPlainEntry()))
	require.NoError(t, group.AddPlain(newVersionedPlainEntry(defaultVersionedFixture("v1"))))
	require.NoError(t, group.AddInstance(newVersionedInstanceEntry(defaultVersionedFixture("v1"))))
	require.NoError(t, group.AddInstance(newVersionedInstanceEntry(defaultVersionedFixture("v2"))))
	require.NoError(t, group.AddOLH(newVersionedOLHEntry("alpha", "v1", false)))
}

func TestEntryGroupClassifierRejectsVersionedObjectWhenOLHReferencesMissingInstance(t *testing.T) {
	t.Parallel()

	group := entrygroup.New("alpha")

	require.NoError(t, group.AddPlain(newVersionedHeadPlainEntry()))
	require.NoError(t, group.AddPlain(newVersionedPlainEntry(defaultVersionedFixture("v1"))))
	require.NoError(t, group.AddInstance(newVersionedInstanceEntry(defaultVersionedFixture("v1"))))
	require.NoError(t, group.AddOLH(newVersionedOLHEntry("alpha", "missing", false)))
}

func TestEntryGroupClassifierRejectsVersionedObjectWhenOLHHasPendingLog(t *testing.T) {
	t.Parallel()

	group := entrygroup.New("alpha")

	require.NoError(t, group.AddPlain(newVersionedHeadPlainEntry()))
	require.NoError(t, group.AddPlain(newVersionedPlainEntry(defaultVersionedFixture("v1"))))
	require.NoError(t, group.AddInstance(newVersionedInstanceEntry(defaultVersionedFixture("v1"))))
	require.NoError(t, group.AddOLH(newVersionedOLHEntry("alpha", "v1", true)))
}

func TestEntryGroupClassifierRejectsUnversionedWhenIdxDiffersFromName(t *testing.T) {
	t.Parallel()

	group := entrygroup.New("alpha")

	require.NoError(t, group.AddPlain(newCustomPlainEntry(
		"alpha-idx",
		"alpha",
		"",
		false,
		1,
		1,
		true,
		sampleMTime(),
		"etag",
		"tag",
		0,
	)))
}

func TestEntryGroupClassifierRejectsUnversionedWhenInstanceIsNotEmpty(t *testing.T) {
	t.Parallel()

	group := entrygroup.New("alpha")

	require.NoError(t, group.AddPlain(newCustomPlainEntry(
		"alpha",
		"alpha",
		"v1",
		false,
		1,
		1,
		true,
		sampleMTime(),
		"etag",
		"tag",
		0,
	)))
}

func TestEntryGroupClassifierRejectsUnversionedWhenVersionPoolIsBelowMinimum(t *testing.T) {
	t.Parallel()

	group := entrygroup.New("alpha")

	require.NoError(t, group.AddPlain(newCustomPlainEntry(
		"alpha",
		"alpha",
		"",
		false,
		0,
		1,
		true,
		sampleMTime(),
		"etag",
		"tag",
		0,
	)))
}

func TestEntryGroupClassifierRejectsUnversionedWhenVersionEpochIsBelowMinimum(t *testing.T) {
	t.Parallel()

	group := entrygroup.New("alpha")

	require.NoError(t, group.AddPlain(newCustomPlainEntry(
		"alpha",
		"alpha",
		"",
		false,
		1,
		0,
		true,
		sampleMTime(),
		"etag",
		"tag",
		0,
	)))
}

func TestEntryGroupClassifierRejectsUnversionedWhenExistsIsFalse(t *testing.T) {
	t.Parallel()

	group := entrygroup.New("alpha")

	require.NoError(t, group.AddPlain(newCustomPlainEntry(
		"alpha",
		"alpha",
		"",
		false,
		1,
		1,
		false,
		sampleMTime(),
		"etag",
		"tag",
		0,
	)))
}

func TestEntryGroupClassifierRejectsUnversionedWhenMTimeIsZero(t *testing.T) {
	t.Parallel()

	group := entrygroup.New("alpha")

	require.NoError(t, group.AddPlain(newCustomPlainEntry(
		"alpha",
		"alpha",
		"",
		false,
		1,
		1,
		true,
		"",
		"etag",
		"tag",
		0,
	)))
}

func TestEntryGroupClassifierRejectsUnversionedWhenETagIsEmpty(t *testing.T) {
	t.Parallel()

	group := entrygroup.New("alpha")

	require.NoError(t, group.AddPlain(newCustomPlainEntry(
		"alpha",
		"alpha",
		"",
		false,
		1,
		1,
		true,
		sampleMTime(),
		"",
		"tag",
		0,
	)))
}

func TestEntryGroupClassifierRejectsUnversionedWhenTagIsEmpty(t *testing.T) {
	t.Parallel()

	group := entrygroup.New("alpha")

	require.NoError(t, group.AddPlain(newCustomPlainEntry(
		"alpha",
		"alpha",
		"",
		false,
		1,
		1,
		true,
		sampleMTime(),
		"etag",
		"",
		0,
	)))
}

func TestEntryGroupClassifierRejectsUnversionedWhenFlagsAreNotZero(t *testing.T) {
	t.Parallel()

	group := entrygroup.New("alpha")

	require.NoError(t, group.AddPlain(newCustomPlainEntry(
		"alpha",
		"alpha",
		"",
		false,
		1,
		1,
		true,
		sampleMTime(),
		"etag",
		"tag",
		1,
	)))
}

func TestEntryGroupRejectsMismatchedPlainName(t *testing.T) {
	t.Parallel()

	group := entrygroup.New("alpha")

	err := group.AddPlain(newPlainEntry("beta", false))
	require.EqualError(
		t,
		err,
		"entry name does not match group name: entry name \"beta\" does not match group name \"alpha\"",
	)
}

func TestEntryGroupRejectsMismatchedInstanceName(t *testing.T) {
	t.Parallel()

	group := entrygroup.New("alpha")

	err := group.AddInstance(newInstanceEntry("beta", false))
	require.EqualError(
		t,
		err,
		"entry name does not match group name: entry name \"beta\" does not match group name \"alpha\"",
	)
}

func TestEntryGroupRejectsMismatchedOLHName(t *testing.T) {
	t.Parallel()

	group := entrygroup.New("alpha")

	err := group.AddOLH(newOLHEntry("beta", false))
	require.EqualError(
		t,
		err,
		"entry name does not match group name: entry name \"beta\" does not match group name \"alpha\"",
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

func newUnversionedPlainEntry() *domain.Plain {
	return newCustomPlainEntry(
		"alpha",
		"alpha",
		"",
		false,
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

func sampleEarlierMTime() string {
	return sampleMTimeAtOffset(-26 * time.Hour)
}

func sampleDeleteMarkerMTime() string {
	return sampleMTimeAtOffset(-23 * time.Hour)
}

func sampleStaleMTime() string {
	return sampleMTimeAtOffset(-8 * 24 * time.Hour)
}

func sampleMTimeAtOffset(offset time.Duration) string {
	now := time.Now().UTC().Truncate(time.Second)

	return now.Add(offset).Format(time.RFC3339Nano)
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
