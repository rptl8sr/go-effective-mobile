insert into tasks(user_id, title, status, created_at, total_duration)
values ($1, $2, 'created', now(), '0 seconds')
returning id;
