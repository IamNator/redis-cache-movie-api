create table if not exists comment
(
    id             uuid      default gen_random_uuid() not null
        constraint comment_pk
            primary key,
    swapi_movie_id int                                 not null,
    message        varchar(600)                        not null,
    ipv4_addr      varchar(20)                         not null,
    created_at     timestamp default current_timestamp not null,
    updated_at     timestamp,
    deleted_at     timestamp
);

comment on column comment.swapi_movie_id is  'id of movie from www.swapi.dev';

comment on column comment.message is 'max length expected is 500; padding = 100 chars';

comment on column comment.ipv4_addr is 'Ip address of the person commenting;
max expected length is 15';

create unique index if not exists comment_id_uindex
    on comment (id);

create index if not exists comment_ipv4_addr_index
    on comment (ipv4_addr);

