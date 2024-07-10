update tasks
set
    status='started',
    started_at=now()
where
    id=$1 and
    user_id=$2