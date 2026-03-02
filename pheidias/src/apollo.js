import { ApolloClient, InMemoryCache, HttpLink } from '@apollo/client/core';
import { provideApolloClient, DefaultApolloClient } from '@vue/apollo-composable';
import fragmentTypes from './constants/fragments.json';

const inMemoryCache = new InMemoryCache({
    possibleTypes: fragmentTypes.__schema.types.reduce((acc, type) => {
        if (type.possibleTypes) {
            acc[type.name] = type.possibleTypes.map(possibleType => possibleType.name);
        }
        return acc;
    }, {})
});

// Define the URL for the GraphQL server
let url = document.location.origin + '/graphql';
if (process.env.NODE_ENV === 'development') {
    url = 'http://localhost:8080/graphql';
}

if (process.env.NODE_ENV === 'local') {
    url = 'http://byzantium.odysseia-greek:4900/graphql';
}

// Create the HTTP link for the Apollo client
const httpLink = new HttpLink({
    uri: url,
    method: 'POST',
    mode: 'cors',
    headers: {
        'Content-Type': 'application/json',
    },
    credentials: 'same-origin',
});

// Create the Apollo client
export const apolloClient = new ApolloClient({
    link: httpLink,
    cache: inMemoryCache,
    connectToDevTools: true,
});

// Provide the Apollo client for use with the Composition API
provideApolloClient(apolloClient);

export default {
    install: (app) => {
        app.provide(DefaultApolloClient, apolloClient);
    }
};
