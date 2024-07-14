select
    id,
    user_id,
    title,
    status,
    coalesce($2::timestamptz, now()) as created_at,
    updated_at at time zone 'utc',
    started_at at time zone 'utc',
    completed_at at time zone 'utc',
    to_char(
        case
            when status = 'stopped'
                then least(completed_at, $3::timestamptz at time zone 'utc')
            when status = 'started'
                then coalesce($3::timestamptz at time zone 'utc', now())
        end
        -
        case
            when $2::timestamptz at time zone 'utc' is null or
                 $2::timestamptz at time zone 'utc' <= started_at
                then started_at
            when $2::timestamptz at time zone 'utc' is not null and
                 $2::timestamptz at time zone 'utc' > started_at
                then $2::timestamptz
        end
    , 'DDD "days" HH24:MI:SS') as total_duration
from tasks
where user_id=$1
  and
    (
        $2::timestamptz is null OR
        $2::timestamptz < completed_at
    ) and
    (
        $3::timestamptz is null OR
        $3::timestamptz > started_at
    )
order by total_duration desc, id