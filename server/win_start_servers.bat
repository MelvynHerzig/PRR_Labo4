for /l %%x in (0, 1, 12) do (
   start cmd /k go run . %%x
)