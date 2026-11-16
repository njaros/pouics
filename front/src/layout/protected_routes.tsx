import { Navigate, Outlet } from "react-router"
import { getToken } from "../api/client"

/**
 * A token control at his simpliest
 */
export default () => {
    const token = getToken()

    if (!token) {
        return <Navigate to={"login"} />
    }

    return <Outlet />
}
