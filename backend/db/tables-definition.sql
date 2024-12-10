/* 1. Subjects Table: Stores topics/categories for games */
CREATE TABLE subjects (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT NOW()
);

/* 2. Games Table: Stores metadata for each game */
CREATE TABLE games (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT,
    subject_id UUID REFERENCES subjects(id) ON DELETE CASCADE,
    difficulty_level INT CHECK (difficulty_level BETWEEN 1 AND 5),
    created_at TIMESTAMP DEFAULT NOW()
);

/* 3. Users Table: Managed by Supabase Auth (optional schema for additional user details) */
CREATE TABLE public.users (
    id uuid not null default gen_random_uuid (),
    email text not null,
    created_at timestamp without time zone null default now(),
    role public.app_role not null default 'viewer'::app_role,
    constraint users_pkey primary key (id),
    constraint users_email_key unique (email)
  );

/* 4. Game States Table: Tracks user-specific progress for active games */
CREATE TABLE game_states (
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    game_id UUID REFERENCES games(id) ON DELETE CASCADE,
    state_data JSONB NOT NULL, /* Stores the current state of the game in JSON format */
    last_updated TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (user_id, game_id)
);

/* 5. Game Results Table: Stores completed game results for users */
CREATE TABLE game_results (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    game_id UUID REFERENCES games(id) ON DELETE CASCADE,
    score INT CHECK (score BETWEEN 0 AND 100),
    completion_time INTERVAL, /* Time taken to complete the game */
    completed_at TIMESTAMP DEFAULT NOW()
);