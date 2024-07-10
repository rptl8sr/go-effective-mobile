select
    id,
    user_id,
    title,
    status,
    created_at,
    updated_at,
    started_at,
    completed_at,
    to_char(total_duration,  'DDD "days" HH24:MI:SS') as total_duration
from tasks
where user_id=$1
order by total_duration desc, id