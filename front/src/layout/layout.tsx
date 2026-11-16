import { Outlet } from "react-router"

export default () => {
    return (
        <div>
            <Header />
            <div>
                <LeftSide />
                <Outlet />
                <RightSide />
            </div>
            <Footer />
        </div>
    )
}