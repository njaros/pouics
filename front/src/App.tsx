import { RouterProvider } from "react-router"
import { UserContext } from "./contexts/user"
import router from "./router/routes"

export default () => {
  return (
    <UserContext value={null}>
      <RouterProvider router={router} />
    </UserContext>
  )
}