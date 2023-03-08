CREATE UNIQUE INDEX IF NOT EXISTS comment_swapi_movie_id_message_ipv4_addr_uindex
    ON busha.public.comment(swapi_movie_id, message, ipv4_addr);
