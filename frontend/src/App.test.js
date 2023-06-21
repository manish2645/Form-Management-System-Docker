import React from 'react';
import { render, fireEvent } from '@testing-library/react';
import App from './App';

describe('App', () => {
  it('renders without errors', () => {
    render(<App />);
  });

  it('displays an error message when required fields are not filled', () => {
    const {getByText} = render(<App />);

    const originalAlert = window.alert;
    window.alert = jest.fn();

    const submitButton = getByText('Submit');
    fireEvent.click(submitButton);

    expect(window.alert).toHaveBeenCalledWith('Please fill in all the required fields.');

    window.alert = originalAlert;
  });

  it('renders "No data available" when leaveData is empty', () => {
    const leaveData = [];

    const { getByText } = render(<App leaveData={leaveData} />);

    expect(getByText('No data available')).toBeInTheDocument();
  });
});
