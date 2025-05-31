// generated code - do not edit
import { type } from "arktype";

export const Todo_Request_Schema = type({
  id: "number",
  title: "string",
});

type Todo_Request = typeof Todo_Request_Schema.infer;
type Todo_Response = {
  id: number;
};

export class RPC_Client {
  constructor(public base_url: string) {}

  async #do_fetch<T>(path: string, args: T) {
    try {
      const result = await fetch(this.base_url + path, {
        method: "POST",
        body: JSON.stringify(args),
      });

      if (!result.ok) {
        throw new Error("fetch result not ok");
      }

      try {
        const data = await result.json();

        return data as T;
      } catch (parse_error) {
        throw parse_error;
      }
    } catch (fetch_error) {
      throw fetch_error;
    }
  }

  add_todo(args: Todo_Request) {
    // validieren?
    return this.#do_fetch<Todo_Response>("/add_todos", args);
  }
}
