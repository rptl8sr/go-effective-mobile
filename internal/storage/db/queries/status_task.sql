select status
from tasks
where
    id=$1 and
    user_id=$2