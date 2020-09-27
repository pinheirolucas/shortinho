create extension if not exists "uuid-ossp";

create table if not exists link (
  created_at timestamp not null default now(),
  updated_at timestamp not null default now(),

  slug       text,
  target_url text not null,
  max_hits   integer,
  active     boolean not null default true,

  primary key (slug)
);

create table if not exists link_hit (
  created_at timestamp not null default now(),
  updated_at timestamp not null default now(),

  id      uuid default uuid_generate_v4(),
  slug    text not null references link (slug),
  deleted boolean not null default false,

  primary key (id)
);
