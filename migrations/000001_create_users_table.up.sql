CREATE TABLE users (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  google_id TEXT UNIQUE NOT NULL,
  primary_email TEXT UNIQUE NOT NULL,
  display_name TEXT,
  profile_image_url TEXT,
  is_archived BOOLEAN DEFAULT FALSE,
  is_banned BOOLEAN DEFAULT FALSE,
  is_muted BOOLEAN DEFAULT FALSE,
  current_mode TEXT,
  is_seller_profile_approved BOOLEAN DEFAULT FALSE,
  is_seller_profile_created BOOLEAN DEFAULT FALSE,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);
