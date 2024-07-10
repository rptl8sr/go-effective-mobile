select
    id,
    user_id,
    title,
    status,
    created_at,
    updated_at,
    started_at,
    completed_at,
case
    when status = 'started'
    then to_char((coalesce($2, now())-started_at),  'DDD "days" HH24:MI:SS')
    else to_char(total_duration,  'DDD "days" HH24:MI:SS')
end as total_duration
from tasks
where id=$1;