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

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Layout />}>
          <Route path="video-add" element={<VideoAdd />} />
          <Route path="video-list" element={<VideoList />} />
        </Route>
      </Routes>
    </BrowserRouter>
  );
}

export default App;
