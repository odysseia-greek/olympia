import gql from 'graphql-tag';

export const CounterTopFive = gql`
  query CounterTopFive {
    counterTopFive {
      topFive {
        lastUsed
        serviceName
        word
        count
      }
    }
  }
`;
