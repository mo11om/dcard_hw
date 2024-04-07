SELECT a.id, a.title, a.start_at, a.end_at
FROM Ads a 
INNER JOIN AdConditions ac ON a.id = ac.ad_id 
INNER JOIN Conditions c ON ac.condition_id = c.id

WHERE
 a.start_at <= NOW() AND a.end_at >= NOW() 
AND 
(a.min_age <= 24 AND a.max_age >= 24)
AND 

((c.type = 'gender' AND c.value = ('F')) 
OR (c.type = 'country' AND c.value = ('TW')) 
OR (c.type = 'platform' AND c.value = ('ios')) )

group by ad_id 
having  COUNT(*)=3
ORDER by end_at ASC



 