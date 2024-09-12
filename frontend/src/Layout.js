import { Outlet, Link } from "react-router-dom";

const Layout = () => {
  return (
    <div className="App">
    <>
      <nav>
        <div>
            <Link to="/video-add">Add Video</Link>
        </div>
        <div>
            <Link to="/video-list">List Videos</Link>
        </div>
      </nav>
      <Outlet />
    </>
    </div>
  )
};

export default Layout;