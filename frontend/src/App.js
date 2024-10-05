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
import { H } from 'highlight.run';

const hProject = process.env.REACT_APP_HIGHLIGHT_PROJECT_ID;

H.init(`${hProject}`, {
	serviceName: "silly-demo-frontend",
	tracingOrigins: true,
	networkRecording: {
		enabled: true,
		recordHeadersAndBody: true,
	},
    enableOtelTracing: true
});
H.identify('viktor@farcic.com', {
  id: 'vfarcic',
  name: 'Viktor Farcic'
})

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
