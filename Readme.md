# Go Backend Svelte Frontend Integration

```go
type AddTodoArgs struct {
  Title string `json:"title" validate:"required,min=3,max=100"`
}

type Todo struct {
  ID    int    `json:"id"`
  Title string `json:"title"`
}
```

Für den Client werden folgende Objekte erzeugt, um

- Validierung
- Type Safe Backend Calls

zu ermöglichen:

- ArkTypes für Args (was an ans Backend geschickt wird)
- typisierte Backend Calls

```Typescript
const AddTodoArgs = type({
  title: "3 <= string <= 100"
});

const Todo = type({
  id: "number",
  title: "string"
});

const fetch = ...

const calls = {
  addTodo: (args: typeof AddTodoArgs) => fetch("/api/todos", {
    method: "POST",
    body: JSON.stringify(args)
  })
};

```

## TODO alle Calls über POST lösen, um es einfach zu halten?
