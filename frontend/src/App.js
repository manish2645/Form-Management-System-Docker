import React, { useState, useEffect } from 'react';
import './App.css'

const App= () => {
  const [name, setName] = useState('');
  const [leaveTypes, setLeaveTypes] = useState([]);
  const [leaveType, setLeaveType] = useState('');
  const [fromDate, setFromDate] = useState('');
  const [toDate, setToDate] = useState('');
  const [team, setTeam] = useState('');
  const [file, setFile] = useState('');
  const [reporter, setReporter] = useState('');

  useEffect(() => {
    fetch('http://localhost:8080/leaveTypes')
      .then((response) => response.json())
      .then((data) => {
        setLeaveTypes(data.leaveTypes);
      })
      .catch((error) => {
        console.error(error);
      });
  }, []);
  
  const [leaveData, setLeaveData] = useState([]);

  useEffect(() => {
    fetch('http://localhost:8080/getleave')
      .then((response) => response.json())
      .then((data) => {
        setLeaveData(data);
      })
      .catch((error) => {
        console.error(error);
      });
  }, []);

  const handleLeaveTypeChange = (e) => {
    setLeaveType(e.target.value);
  };

  const handleFormSubmit = (e) => {
    e.preventDefault();
    
    if (!name || !leaveType || !fromDate || !toDate || !team || !reporter) {
      alert('Please fill in all the required fields.');
      return;
    }

    const formData = new FormData();
    formData.append('name', name);
    formData.append('leaveType', leaveType);
    formData.append('fromDate', fromDate);
    formData.append('toDate', toDate);
    formData.append('team', team);
    formData.append('reporter', reporter);
  
    if (leaveType === 'Sick Leave' && file) {
      formData.append('file', file);
    }
  
    fetch('http://localhost:8080/postleave', {
      method: 'POST',
      body: formData,
    })
      .then((response) => {
        if (response.ok) {
          window.alert('Leave application submitted successfully');
        } else {
          throw new Error('Failed to submit leave application');
        }
      })
      .catch((error) => {
        console.error(error);
        window.alert('Failed! Try Again?');
      });

      setName('');
      setLeaveType('');
      setFromDate('');
      setToDate('');
      setFile('');
      setTeam('');
      setReporter('');

  };

  const openAttachment = (filePath) => {
    window.open(`http://localhost:8080/file/${filePath}`);
  };

  return (
    <div className='outerContainer'>
      <div className="container">
        <div className="box">
          <h1>Leave Form</h1>
          <form onSubmit={handleFormSubmit}>
            <label htmlFor="name">Name:</label>
            <input
              type="text"
              id="name"
              value={name}
              onChange={(e) => setName(e.target.value)}
              required />
            <div>
              <label>
                Leave Type:
                {Array.isArray(leaveTypes) && leaveTypes.map((type) => (
                  <label key={type}>
                    <input
                      type="radio"
                      name="leaveType"
                      value={type}
                      checked={leaveType === type}
                      onChange={handleLeaveTypeChange}
                      required />
                    {type}
                  </label>
                ))}
              </label>
            </div>
            <br />
            <label htmlFor="fromDate">From:</label>
            <input
              type="date"
              id="fromDate"
              value={fromDate}
              onChange={(e) => setFromDate(e.target.value)}
              required />
            <br />
            <label htmlFor="toDate">To:</label>
            <input
              type="date"
              id="toDate"
              value={toDate}
              onChange={(e) => setToDate(e.target.value)}
              required />

            <div>
              <br />
              <label>
                Team Names:
                <input
                  type="radio"
                  name="team"
                  value="CloudOps"
                  checked={team === 'CloudOps'}
                  onChange={(e) => setTeam(e.target.value)}
                  required />
                CloudOps
              </label>

              <label>
                <input
                  type="radio"
                  name="team"
                  value="DataOps"
                  checked={team === 'DataOps'}
                  onChange={(e) => setTeam(e.target.value)}
                  required />
                DataOps
              </label>

              <label>
                <input
                  type="radio"
                  name="team"
                  value="AnalyticOps"
                  checked={team === 'AnalyticOps'}
                  onChange={(e) => setTeam(e.target.value)}
                  required />
                AnalyticOps
              </label>
            </div>
            <br />
            {leaveType === 'Sick Leave' && (
              <div>
                <label htmlFor="file">File Upload (max: 15mb, pdf/png):</label>
                <input
                  type="file"
                  id="file"
                  onChange={(e) => setFile(e.target.files[0])}
                  accept=".pdf,.png"
                  required />
              </div>
            )}
            <br />
            <label htmlFor="reporter">Reporter:</label>
            <select
              id="reporter"
              value={reporter}
              onChange={(e) => setReporter(e.target.value)}
              required
            >
              <option value="">Select Reporter</option>
              <option value="Surya Kant">Surya Kant</option>
              <option value="Pradeep Kumar Bharti">Pradeep Kumar Bharti</option>
              <option value="Nitin Aggarwal">Nitin Aggarwal</option>
              <option value="Avinashi Sharma">Avinashi Sharma</option>
            </select>
            <br />
            <button type="submit">Submit</button>
          </form>
          <br />
        </div>
      </div>

      <div className='container'>
        <div className='box'>
          <h1>Applied Leaves</h1>
          <table className="leave-table">
            <thead>
              <tr>
                <th>Serial No</th>
                <th>Name</th>
                <th>Leave Type</th>
                <th>From</th>
                <th>To</th>
                <th>Team</th>
                <th>Sick Leave Attachment</th>
                <th>Reporter</th>
              </tr>
            </thead>
            <tbody>
              {Array.isArray(leaveData) && leaveData.length > 0 ? (
                leaveData.map((leave) => (
                  <tr key={leave.id}>
                    <td>{leave.id}</td>
                    <td>{leave.name}</td>
                    <td>{leave.leaveType}</td>
                    <td>{leave.fromDate}</td>
                    <td>{leave.toDate}</td>
                    <td>{leave.team}</td>
                    <td>
                      <p className="attachment-link" onClick={() => openAttachment(leave.filePath)}><u>{leave.filePath}</u></p>
                    </td>
                    <td>{leave.reporter}</td>
                  </tr>
                ))
              ) : (
                <tr>
                  <td colSpan="7">No data available</td>
                </tr>
              )}
            </tbody>
          </table>
        </div>
      </div>
    </div>

  );
};

export default App;
