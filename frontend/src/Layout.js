import { Outlet, Link } from "react-router-dom";
import Root from './Root';

const url = process.env.BACKEND_URL;

const Layout = () => {
  return (
    <div className="App">
    <>
      <Root url={url} />
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