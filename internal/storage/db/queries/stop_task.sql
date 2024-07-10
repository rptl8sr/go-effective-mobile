update tasks
set
    status='stopped',
    completed_at=now(),
    total_duration=total_duration + (now()-started_at)
where id=$1
returning
    concat(
        to_char(total_duration,  'DDD "days" HH24:MI:SS')
    );