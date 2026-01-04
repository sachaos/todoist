# Todoist API Migration Plan: v9 → v1

## Executive Summary

This document outlines the minimal migration strategy from v9 Sync API to v1 API.

**Current State**: `https://todoist.com/API/v9/`
**Target State**: `https://api.todoist.com/api/v1/`

**Key Finding**: The codebase is already largely compatible with v1 API. Only 3 endpoint changes required.

**Decisions**:
- **Cache Strategy**: Invalidate on error (defensive), no proactive migration needed since IDs are already strings
- **Architecture**: Keep existing Sync API patterns, no REST API introduction
- **Scope**: Minimal changes only - strictly what migration guide requires

---

## Migration Scope Analysis

### Already Compatible (No Changes Needed)

| Aspect | Current State | v1 Requirement | Status |
|--------|--------------|----------------|--------|
| ID types | `string` | `string` | ✓ Compatible |
| Map types | `map[string]*Item` etc. | `map[string]*...` | ✓ Compatible |
| Labels field | `[]string` (label names) with json `"labels"` | `labels: list[str]` | ✓ Compatible |
| Sync command types | `item_add`, `item_close`, etc. | Same | ✓ Compatible |
| Content-Type | `application/x-www-form-urlencoded` | Same for Sync API | ✓ Compatible |
| Sync endpoint path | `/sync` | `/sync` | ✓ Compatible |
| Bearer token auth | Yes | Yes | ✓ Compatible |

### Required Changes (3 Total)

#### 1. Base URL Change
**File**: `lib/main.go:15`
```go
// Current
const Server = "https://todoist.com/API/v9/"

// Target
const Server = "https://api.todoist.com/api/v1/"
```

#### 2. Quick Add Endpoint
**File**: `lib/todoist.go:106`
```go
// Current
return c.doApi(ctx, http.MethodPost, "quick/add", values, &r)

// Target
return c.doApi(ctx, http.MethodPost, "tasks/quick_add", values, &r)
```

#### 3. Completed Tasks Endpoint
**File**: `lib/completed.go:15`
```go
// Current
return c.doApi(ctx, http.MethodPost, "completed/get_all", url.Values{}, &r)

// Target (note: may need to verify exact endpoint and response format)
return c.doApi(ctx, http.MethodGet, "tasks/completed_by_completion_date", url.Values{}, &r)
```

**Note**: The completed tasks endpoint may require additional investigation:
- v1 uses `GET` method (vs `POST` in v9)
- Response format may differ
- May need pagination handling

---

## Implementation Checklist

### Phase 1: Core URL Change
- [ ] Update `lib/main.go:15` - Change base URL constant
- [ ] Test basic sync operation with v1 endpoint

### Phase 2: Deprecated Endpoint Updates
- [ ] Update `lib/todoist.go:106` - Quick add endpoint
- [ ] Update `lib/completed.go:15` - Completed tasks endpoint
- [ ] Verify response format for completed tasks (may need struct updates)

### Phase 3: Defensive Cache Handling
- [ ] Add graceful error handling in `cache.go` if unmarshal fails
- [ ] Log message about cache invalidation
- [ ] Force re-sync on parse error

### Phase 4: Testing
- [ ] Test `sync` command
- [ ] Test `add` command
- [ ] Test `modify` command
- [ ] Test `close` command
- [ ] Test `delete` command
- [ ] Test `quick` command
- [ ] Test `completed-list` command
- [ ] Test filters work correctly

---

## Risk Assessment

| Risk | Likelihood | Impact | Mitigation |
|------|------------|--------|------------|
| Base URL change breaks sync | Low | High | Test immediately after change |
| Quick add response format differs | Medium | Low | Verify response parsing |
| Completed endpoint response differs | High | Medium | Investigate v1 response format |
| Cache incompatibility | Low | Low | Graceful cache invalidation |

---

## References

- [Todoist API v1 Documentation](https://developer.todoist.com/api/v1)
- [Migration Guide](https://developer.todoist.com/api/v1#tag/Migrating-from-v9)
- Base URL: `https://api.todoist.com/api/v1/`
- Sync endpoint: `https://api.todoist.com/api/v1/sync`
