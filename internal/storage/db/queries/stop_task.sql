update tasks
set
    status='stopped',
    completed_at=now(),
    total_duration=now()-started_at
where id=$1
returning
concat(
    extract(days from total_duration), ' days',
    to_char(now()-started_at, ' HH24:MI:SS')
);
