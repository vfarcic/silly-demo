import './App.css';
import React from 'react';
import {
  BrowserRouter,
  Routes,
  Route,
} from "react-router-dom";
import Layout from './Layout';
import VideoList from './VideoList';
import VideoAdd from './VideoAdd';

const url = process.env.BACKEND_URL;

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Layout />}>
          <Route path="video-add" element={<VideoAdd url={url} />} />
          <Route path="video-list" element={<VideoList url={url} />} />
        </Route>
      </Routes>
    </BrowserRouter>
  );
}

export default App;
