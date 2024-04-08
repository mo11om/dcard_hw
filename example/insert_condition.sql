-- Replace `<ad_id>` with the actual ad ID (retrieved in step 2 or your existing ad)
INSERT INTO AdConditions (ad_id, condition_id)
SELECT <ad_id>, c.id
FROM Conditions c
WHERE  
  OR c.type = 'gender' AND c.value IN ('M', 'F')  -- Target both genders
  OR c.type = 'country' AND c.value IN ('TW', 'JP')  -- Target both countries (Taiwan & Japan)
  OR c.type = 'platform' AND c.value IN ('android', 'ios', 'web');  -- Target all platforms
