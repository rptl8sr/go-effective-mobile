update tasks
set
    status='stopped',
    completed_at=now(),
    total_duration=now()-started_at
where id=$1
returning total_duration;
