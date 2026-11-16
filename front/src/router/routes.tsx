import { createBrowserRouter } from "react-router";
import Layout from "../layout/layout";

export default createBrowserRouter([
    {
        element: <ProtectedRoutes />,
        children: [
            {
                path: "/",
                element: <Layout />,
                children: [
                    {
                        path: "/",
                        element: <Home />
                    },
                    {
                        path: "/game",
                        element: <Game />
                    }
                ]
            },
        ]
    },
    {
        path: "/login",
        element: <Login />
    },
    {
        path: "/register",
        element: <Register />
    },
    {
        path: "*",
        element: <NotFound />
    }
])