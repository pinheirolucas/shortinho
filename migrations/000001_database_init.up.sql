create extension if not exists "uuid-ossp";

create table if not exists link (
  created_at timestamp not null,
  updated_at timestamp not null,

  slug       text,
  target_url text not null,
  max_hits   integer,

  primary key (slug)
);

create table if not exists link_hit (
  created_at timestamp not null,
  updated_at timestamp not null,

  id   uuid default uuid_generate_v4(),
  slug text not null references link (slug),

  primary key (id)
);
