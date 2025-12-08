CREATE TABLE system_check_log (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  test_label TEXT,
  created_at TIMESTAMP DEFAULT NOW()
);
