for /l %%x in (0, 1, 7) do (
   start cmd /k go run . %%x
)