CREATE TABLE events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    userid UUID NOT NULL REFERENCES users(id) ON DELETE SET NULL,  -- optional but safe FK

    event_type TEXT NOT NULL,
    event_payload JSONB NOT NULL,

    dispatched_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);


-- For dispatching outbox patterns
CREATE INDEX idx_events_dispatched_at ON events(dispatched_at);

-- For filtering event types
CREATE INDEX idx_events_event_type ON events(event_type);

-- For actor-based event history
CREATE INDEX idx_events_userid ON events(userid);